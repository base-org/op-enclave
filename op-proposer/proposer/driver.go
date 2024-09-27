package proposer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum-optimism/optimism/op-proposer/metrics"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/mdehoog/op-nitro/bindings"
	"github.com/mdehoog/op-nitro/enclave"
)

var (
	ErrProposerNotRunning = errors.New("proposer is not running")
)

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

	prover *prover
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

	prover, err := newProver(cCtx, setup.L1Client, setup.L2Client, setup.RollupClient, setup.EnclaveClient)
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

// loop is responsible for creating & submitting the next outputs
// The loop regularly polls the L2 chain to infer whether to make the next proposal.
func (l *L2OutputSubmitter) loop() {
	defer l.wg.Done()
	defer l.Log.Info("loop returning")
	ctx := l.ctx
	ticker := time.NewTicker(l.Cfg.PollInterval)
	defer ticker.Stop()
	var lastProposal *proposal
	for {
		select {
		case <-ticker.C:
			// prioritize quit signal
			select {
			case <-l.done:
				return
			default:
			}

			// A note on retrying: the outer ticker already runs on a short
			// poll interval, which has a default value of 6 seconds. So no
			// retry logic is needed around output fetching here.
			proposal, shouldPropose, err := l.generateNextProposal(ctx, lastProposal)
			lastProposal = proposal
			if err != nil {
				l.Log.Warn("Error getting output", "err", err)
				continue
			} else if !shouldPropose {
				// debug logging already in Fetch(DGF|L2OO)Output
				continue
			}

			l.proposeOutput(ctx, proposal)
		case <-l.done:
			return
		}
	}
}

func (l *L2OutputSubmitter) generateNextProposal(ctx context.Context, lastProposal *proposal) (*proposal, bool, error) {
	proposed, err := l.ooContract.LatestL2Output(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return nil, false, fmt.Errorf("failed to get latest proposed block number from Oracle: %w", err)
	}
	proposedBlockNumber := proposed.L2BlockNumber.Uint64()

	// purge reorged blocks
	if lastProposal != nil {
		lastProposalBlockRef := lastProposal.blockRef
		proposedHeader, err := l.L2Client.HeaderByNumber(ctx, new(big.Int).SetUint64(lastProposalBlockRef.Number))
		if err != nil {
			return nil, false, fmt.Errorf("failed to get header for block %d: %w", lastProposalBlockRef.Number, err)
		}
		if lastProposalBlockRef.Hash != proposedHeader.Hash() {
			l.Log.Warn("Last proposal block hash does not match the L2 block hash, possible reorg",
				"last_proposal", lastProposalBlockRef.Hash, "l2_block", proposedHeader.Hash())
			// TODO rather than clearing all aggregated proposals, store snapshots and binary search back to the common ancestor
			lastProposal = nil
		} else {
			proposedBlockNumber = lastProposalBlockRef.Number
		}
	}

	// generate new proposals up to the latest block
	syncStatus, err := l.RollupClient.SyncStatus(ctx)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get sync status from Rollup: %w", err)
	}
	latestBlockNumber := syncStatus.FinalizedL2.Number
	if l.Cfg.AllowNonFinalized {
		latestBlockNumber = syncStatus.SafeL2.Number
	}

	// TODO implement proposal array limit (aggregate in blocks)
	// TODO implement a pool of go-routines for parallel proof generation
	var proposals []*proposal
	if lastProposal != nil {
		proposals = append(proposals, lastProposal)
	}
	shouldPropose := l.Cfg.MinProposalInterval > 0 && latestBlockNumber-proposedBlockNumber > l.Cfg.MinProposalInterval
	for i := proposedBlockNumber + 1; i <= latestBlockNumber; i++ {
		proposal, anyWithdrawals, err := l.prover.Generate(ctx, i)
		if err != nil {
			return nil, false, fmt.Errorf("failed to generate proof for block %d: %w", i, err)
		}
		proposals = append(proposals, proposal)
		shouldPropose = shouldPropose || anyWithdrawals
	}

	lastProposal, err = l.prover.Aggregate(ctx, proposed.OutputRoot, proposals)
	if err != nil {
		return nil, false, fmt.Errorf("failed to aggregate proofs: %w", err)
	}
	return lastProposal, shouldPropose, nil
}

func (l *L2OutputSubmitter) proposeOutput(ctx context.Context, proposal *proposal) {
	cCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	if err := l.sendTransaction(cCtx, proposal); err != nil {
		l.Log.Error("Failed to send proposal transaction",
			"err", err,
			"block", proposal.blockRef)
		return
	}
	l.Metr.RecordL2BlocksProposed(proposal.blockRef)

	// TODO purge witnesses here
}

// sendTransaction creates & sends transactions through the underlying transaction manager.
func (l *L2OutputSubmitter) sendTransaction(ctx context.Context, proposal *proposal) error {
	l.Log.Info("Proposing output root", "output", proposal.output.OutputRoot, "block", proposal.blockRef)
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
func (l *L2OutputSubmitter) ProposeL2OutputTxData(proposal *proposal) ([]byte, error) {
	return proposeL2OutputTxData(l.ooABI, proposal)
}

// proposeL2OutputTxData creates the transaction data for the ProposeL2Output function
func proposeL2OutputTxData(abi *abi.ABI, proposal *proposal) ([]byte, error) {
	return abi.Pack(
		"proposeL2Output",
		proposal.output.OutputRoot,
		new(big.Int).SetUint64(proposal.blockRef.Number),
		new(big.Int).SetUint64(proposal.blockRef.L1Origin.Number),
		[]byte(proposal.output.Signature),
	)
}
