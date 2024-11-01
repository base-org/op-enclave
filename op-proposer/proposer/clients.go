package proposer

import (
	"context"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/sources/caching"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

type L1Client interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	BlockReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error)
	CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	Close()
}

type L2Client interface {
	ChainConfig(ctx context.Context) (*params.ChainConfig, error)
	GetProof(ctx context.Context, address common.Address, hash common.Hash) (*eth.AccountResult, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	ExecutionWitness(ctx context.Context, hash common.Hash) ([]byte, error)
	Close()
}

type Client interface {
	L1Client
	L2Client
}

type RollupClient interface {
	RollupConfig(ctx context.Context) (*rollup.Config, error)
	SyncStatus(ctx context.Context) (*eth.SyncStatus, error)
}

type ethClient struct {
	client        *ethclient.Client
	blocksCache   *caching.LRUCache[common.Hash, *types.Block]
	headersCache  *caching.LRUCache[common.Hash, *types.Header]
	receiptsCache *caching.LRUCache[common.Hash, types.Receipts]
	proofsCache   *caching.LRUCache[[common.AddressLength + common.HashLength]byte, *eth.AccountResult]
}

func NewClient(client *ethclient.Client, metrics caching.Metrics) Client {
	return newClient(client, metrics)
}

func newClient(client *ethclient.Client, metrics caching.Metrics) *ethClient {
	cacheSize := 1000
	return &ethClient{
		client:        client,
		blocksCache:   caching.NewLRUCache[common.Hash, *types.Block](metrics, "blocks", cacheSize),
		headersCache:  caching.NewLRUCache[common.Hash, *types.Header](metrics, "headers", cacheSize),
		receiptsCache: caching.NewLRUCache[common.Hash, types.Receipts](metrics, "receipts", cacheSize),
		proofsCache:   caching.NewLRUCache[[common.AddressLength + common.HashLength]byte, *eth.AccountResult](metrics, "proofs", cacheSize),
	}
}

func (e *ethClient) ChainConfig(ctx context.Context) (*params.ChainConfig, error) {
	var config params.ChainConfig
	err := e.client.Client().CallContext(ctx, &config, "debug_chainConfig")
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (e *ethClient) BlockNumber(ctx context.Context) (uint64, error) {
	return e.client.BlockNumber(ctx)
}

func (e *ethClient) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	if header, ok := e.headersCache.Get(hash); ok {
		return header, nil
	}
	header, err := e.client.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	e.headersCache.Add(hash, header)
	return header, nil
}

func (e *ethClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	header, err := e.client.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	e.headersCache.Add(header.Hash(), header)
	return header, nil
}

func (e *ethClient) BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error) {
	block, err := e.client.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	e.blocksCache.Add(block.Hash(), block)
	e.headersCache.Add(block.Hash(), block.Header())
	return block, nil
}

func (e *ethClient) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	if block, ok := e.blocksCache.Get(hash); ok {
		return block, nil
	}
	block, err := e.client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	e.blocksCache.Add(block.Hash(), block)
	e.headersCache.Add(block.Hash(), block.Header())
	return block, nil
}

func (e *ethClient) BlockReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	if receipts, ok := e.receiptsCache.Get(hash); ok {
		return receipts, nil
	}
	receipts, err := e.client.BlockReceipts(ctx, rpc.BlockNumberOrHash{BlockHash: &hash})
	if err != nil {
		return nil, err
	}
	e.receiptsCache.Add(hash, receipts)
	return receipts, nil
}

func (e *ethClient) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return e.client.CodeAt(ctx, contract, blockNumber)
}

func (e *ethClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	return e.client.CallContract(ctx, call, blockNumber)
}

func (e *ethClient) GetProof(ctx context.Context, address common.Address, hash common.Hash) (*eth.AccountResult, error) {
	key := [common.AddressLength + common.HashLength]byte{}
	copy(key[:common.AddressLength], address[:])
	copy(key[common.AddressLength:], hash[:])
	if proof, ok := e.proofsCache.Get(key); ok {
		return proof, nil
	}
	var proof *eth.AccountResult
	err := e.client.Client().CallContext(ctx, &proof, "eth_getProof", address, []common.Hash{}, hash)
	if err != nil {
		return nil, err
	}
	if proof == nil {
		return nil, ethereum.NotFound
	}
	e.proofsCache.Add(key, proof)
	return proof, nil
}

func (e *ethClient) ExecutionWitness(ctx context.Context, hash common.Hash) ([]byte, error) {
	var buf hexutil.Bytes
	err := e.client.Client().CallContext(ctx, &buf, "debug_executionWitness", hash)
	return buf, err
}

func (e *ethClient) Close() {
	e.client.Close()
}

type rollupClient struct {
	client       *rpc.Client
	witnessCache *caching.LRUCache[common.Hash, []byte]
}

func NewRollupClient(client *rpc.Client, metrics caching.Metrics) RollupClient {
	cacheSize := 1000
	return &rollupClient{
		client:       client,
		witnessCache: caching.NewLRUCache[common.Hash, []byte](metrics, "witnesses", cacheSize),
	}
}

func (w *rollupClient) RollupConfig(ctx context.Context) (*rollup.Config, error) {
	var config rollup.Config
	err := w.client.CallContext(ctx, &config, "optimism_rollupConfig")
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (w *rollupClient) SyncStatus(ctx context.Context) (*eth.SyncStatus, error) {
	var status eth.SyncStatus
	err := w.client.CallContext(ctx, &status, "optimism_syncStatus")
	if err != nil {
		return nil, err
	}
	return &status, nil
}
