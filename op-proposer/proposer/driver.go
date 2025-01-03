package proposer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/base-org/op-enclave/bindings"
	"github.com/base-org/op-enclave/op-enclave/enclave"
	"github.com/ethereum-optimism/optimism/op-proposer/metrics"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
)

const aggregateBatchSize = 1000

var (
	ErrProposerNotRunning = errors.New("proposer is not running")
)

func init() {
	if aggregateBatchSize <= 1 {
		panic("aggregateBatchSize must be greater than 1")
	}
}

type OOContract interface {
	Version(*bind.CallOpts) (string, error)
	LatestL2Output(opts *bind.CallOpts) (bindings.TypesOutputProposal, error)
}

type DriverSetup struct {
	Log           log.Logger
	Metr          metrics.Metricer
	Cfg           ProposerConfig
	Txmgr         txmgr.TxManager
	L1Client      L1Client
	L2Client      L2Client
	RollupClient  RollupClient
	EnclaveClient enclave.RPC
}

// L2OutputSubmitter is responsible for proposing outputs
type L2OutputSubmitter struct {
	DriverSetup

	wg   sync.WaitGroup
	done chan struct{}

	ctx    context.Context
	cancel context.CancelFunc

	mutex   sync.Mutex
	running bool

	ooContract OOContract
	ooABI      *abi.ABI

	prover  *Prover
	pending []*Proposal
}

// NewL2OutputSubmitter creates a new L2 Output Submitter
func NewL2OutputSubmitter(setup DriverSetup) (_ *L2OutputSubmitter, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	// The above context is long-lived, and passed to the `L2OutputSubmitter` instance. This context is closed by
	// `StopL2OutputSubmitting`, but if this function returns an error or panics, we want to ensure that the context
	// doesn't leak.
	defer func() {
		if err != nil || recover() != nil {
			cancel()
		}
	}()

	return newL2OOSubmitter(ctx, cancel, setup)
}

func newL2OOSubmitter(ctx context.Context, cancel context.CancelFunc, setup DriverSetup) (*L2OutputSubmitter, error) {
	ooContract, err := bindings.NewOutputOracleCaller(*setup.Cfg.L2OutputOracleAddr, setup.L1Client)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create L2OO at address %s: %w", setup.Cfg.L2OutputOracleAddr, err)
	}

	cCtx, cCancel := context.WithTimeout(ctx, setup.Cfg.NetworkTimeout)
	defer cCancel()
	version, err := ooContract.Version(&bind.CallOpts{Context: cCtx})
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("Connected to L2OutputOracle", "address", setup.Cfg.L2OutputOracleAddr, "version", version)

	parsed, err := bindings.OutputOracleMetaData.GetAbi()
	if err != nil {
		cancel()
		return nil, err
	}

	prover, err := NewProver(cCtx, setup.L1Client, setup.L2Client, setup.RollupClient, setup.EnclaveClient)
	if err != nil {
		cancel()
		return nil, err
	}

	return &L2OutputSubmitter{
		DriverSetup: setup,
		done:        make(chan struct{}),
		ctx:         ctx,
		cancel:      cancel,

		ooContract: ooContract,
		ooABI:      parsed,
		prover:     prover,
	}, nil
}

func (l *L2OutputSubmitter) StartL2OutputSubmitting() error {
	l.Log.Info("Starting Proposer")

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.running {
		return errors.New("proposer is already running")
	}
	l.running = true

	l.wg.Add(1)
	go l.loop()

	l.Log.Info("Proposer started")
	return nil
}

func (l *L2OutputSubmitter) StopL2OutputSubmittingIfRunning() error {
	err := l.StopL2OutputSubmitting()
	if errors.Is(err, ErrProposerNotRunning) {
		return nil
	}
	return err
}

func (l *L2OutputSubmitter) StopL2OutputSubmitting() error {
	l.Log.Info("Stopping Proposer")

	l.mutex.Lock()
	defer l.mutex.Unlock()

	if !l.running {
		return ErrProposerNotRunning
	}
	l.running = false

	l.cancel()
	close(l.done)
	l.wg.Wait()

	l.Log.Info("Proposer stopped")
	return nil
}

func (l *L2OutputSubmitter) latestSafeBlock(ctx context.Context) (eth.L2BlockRef, error) {
	syncStatus, err := l.RollupClient.SyncStatus(ctx)
	if err != nil {
		return eth.L2BlockRef{}, fmt.Errorf("failed to get sync status from Rollup: %w", err)
	}
	batched := syncStatus.FinalizedL2
	if l.Cfg.AllowNonFinalized {
		batched = syncStatus.SafeL2
	}
	return batched, nil
}

func l2BlockRefToBlockID(ref eth.L2BlockRef) eth.BlockID {
	return eth.BlockID{
		Number: ref.Number,
		Hash:   ref.Hash,
	}
}

// loop is responsible for creating & submitting the next outputs
// The loop regularly polls the L2 chain to infer whether to make the next proposal.
func (l *L2OutputSubmitter) loop() {
	defer l.wg.Done()
	defer l.Log.Info("loop returning")
	ctx := l.ctx
	ticker := time.NewTicker(l.Cfg.PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// prioritize quit signal
			select {
			case <-l.done:
				return
			default:
			}

			latestOutput, err := l.ooContract.LatestL2Output(&bind.CallOpts{Context: ctx})
			if err != nil {
				log.Warn("Failed to get latest proposed block number from Oracle", "err", err)
				continue
			}

			if err = l.generateOutputs(ctx, latestOutput); err != nil {
				l.Log.Warn("Error generating output", "err", err)
				continue
			}

			proposal, shouldPropose, err := l.nextOutput(ctx, latestOutput)
			if err != nil {
				l.Log.Warn("Error getting output", "err", err)
				continue
			} else if !shouldPropose {
				continue
			}

			l.proposeOutput(ctx, proposal)
		case <-l.done:
			return
		}
	}
}

func (l *L2OutputSubmitter) IsRunning() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.running
}

func (l *L2OutputSubmitter) generateOutputs(ctx context.Context, latestOutput bindings.TypesOutputProposal) error {
	latestOutputNumber := latestOutput.L2BlockNumber.Uint64()

	// clear out already submitted outputs
	for len(l.pending) > 0 && l.pending[0].From.Number-1 < latestOutputNumber {
		l.pending = l.pending[1:]
	}

	if len(l.pending) > 0 {
		if l.pending[0].From.Number-1 != latestOutputNumber {
			l.Log.Warn("Pending outputs are not contiguous with the latest output",
				"latest", latestOutputNumber,
				"pending", l.pending[0].From.Number-1)
			l.pending = nil
		} else {
			latestOutputNumber = l.pending[len(l.pending)-1].To.Number
		}
	}

	// calculate `aggregateBatchSize` proofs at once, which are then aggregated in `nextOutput`
	for i := uint64(0); i < aggregateBatchSize; i++ {
		number := i + latestOutputNumber + 1
		block, err := l.L2Client.BlockByNumber(ctx, new(big.Int).SetUint64(number))
		if errors.Is(err, ethereum.NotFound) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to get block %d: %w", number, err)
		}

		proposal, err := l.prover.Generate(ctx, block)
		if err != nil {
			return fmt.Errorf("failed to generate proof for block %d: %w", number, err)
		}

		l.Log.Info("Generated proof for block",
			"block", l2BlockRefToBlockID(proposal.To), "l1Origin", proposal.To.L1Origin,
			"withdrawals", proposal.Withdrawals, "output", proposal.Output.OutputRoot.String())
		l.pending = append(l.pending, proposal)
	}

	return nil
}

func (l *L2OutputSubmitter) nextOutput(ctx context.Context, latestOutput bindings.TypesOutputProposal) (*Proposal, bool, error) {
	// aggregate proposals up to the latest safe block
	latestSafe, err := l.latestSafeBlock(ctx)
	if err != nil {
		return nil, false, err
	}

	count := 0
	for count < len(l.pending) && l.pending[count].To.Number <= latestSafe.Number {
		count++
	}
	if count <= 0 {
		return nil, false, nil
	}

	for count > 1 {
		batchLength := min(count, aggregateBatchSize)
		batch := l.pending[:batchLength]
		aggregated, err := l.prover.Aggregate(ctx, latestOutput.OutputRoot, batch)
		if err != nil {
			var rpcError rpc.Error
			if errors.As(err, &rpcError) {
				// if we received an explicit error from the enclave (like "invalid signer"), clear the pending proofs
				l.Log.Warn("Non-recoverable error aggregating proofs", "err", err)
				l.pending = nil
			}
			return nil, false, err
		}
		l.pending = append([]*Proposal{aggregated}, l.pending[batchLength:]...)
		count -= batchLength - 1
		l.Log.Info("Aggregated proofs",
			"output", aggregated.Output.OutputRoot.String(), "blocks", batchLength, "remaining", count-1,
			"withdrawals", aggregated.Withdrawals, "from", aggregated.From.Number, "to", aggregated.To.Number)
	}
	proposal := l.pending[0]

	if proposal.To.Number < latestSafe.Number {
		l.Log.Info("Aggregated output is not the latest safe block, waiting for more proofs",
			"aggregated", l2BlockRefToBlockID(proposal.To), "latestSafe", latestSafe)
		return nil, false, nil
	}

	if proposal.To.Hash != latestSafe.Hash {
		l.Log.Warn("Aggregated output does not match the latest batched block, possible reorg",
			"aggregated", l2BlockRefToBlockID(proposal.To), "latestSafe", latestSafe)
		l.pending = nil
		return nil, false, nil
	}

	shouldPropose := proposal.Withdrawals ||
		(l.Cfg.MinProposalInterval > 0 &&
			latestSafe.Number-latestOutput.L2BlockNumber.Uint64() > l.Cfg.MinProposalInterval)

	if shouldPropose {
		latestL1Number, err := l.L1Client.BlockNumber(l.ctx)
		if err != nil {
			log.Warn("Failed to get latest block header", "err", err)
			shouldPropose = false
		} else if proposal.To.L1Origin.Number <= latestL1Number-(256-10) {
			// only submit onchain if within the blockhash window - 10
			log.Warn("Not submitting proposal, block is too old", "l1Origin", proposal.To.L1Origin.Number, "l1Latest", latestL1Number)
			shouldPropose = false
		}
	}

	return proposal, shouldPropose, nil
}

func (l *L2OutputSubmitter) proposeOutput(ctx context.Context, proposal *Proposal) {
	cCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	if err := l.sendTransaction(cCtx, proposal); err != nil {
		l.Log.Error("Failed to send proposal transaction",
			"err", err,
			"block", l2BlockRefToBlockID(proposal.To))
		return
	}
	l.Metr.RecordL2BlocksProposed(proposal.To)
}

// sendTransaction creates & sends transactions through the underlying transaction manager.
func (l *L2OutputSubmitter) sendTransaction(ctx context.Context, proposal *Proposal) error {
	l.Log.Info("Proposing output root", "output", proposal.Output.OutputRoot, "block", l2BlockRefToBlockID(proposal.To))
	data, err := l.ProposeL2OutputTxData(proposal)
	if err != nil {
		return err
	}
	receipt, err := l.Txmgr.Send(ctx, txmgr.TxCandidate{
		TxData:   data,
		To:       l.Cfg.L2OutputOracleAddr,
		GasLimit: 0,
	})
	if err != nil {
		return err
	}

	if receipt.Status == types.ReceiptStatusFailed {
		l.Log.Error("Proposer tx successfully published but reverted", "tx_hash", receipt.TxHash)
	} else {
		l.Log.Info("Proposer tx successfully published", "tx_hash", receipt.TxHash)
	}
	return nil
}

// ProposeL2OutputTxData creates the transaction data for the ProposeL2Output function
func (l *L2OutputSubmitter) ProposeL2OutputTxData(proposal *Proposal) ([]byte, error) {
	return proposeL2OutputTxData(l.ooABI, proposal)
}

// proposeL2OutputTxData creates the transaction data for the ProposeL2Output function
func proposeL2OutputTxData(abi *abi.ABI, proposal *Proposal) ([]byte, error) {
	sig := make([]byte, len(proposal.Output.Signature))
	copy(sig, proposal.Output.Signature)
	sig[64] += 27
	return abi.Pack(
		"proposeL2Output",
		proposal.Output.OutputRoot,
		new(big.Int).SetUint64(proposal.To.Number),
		new(big.Int).SetUint64(proposal.To.L1Origin.Number),
		sig,
	)
}
