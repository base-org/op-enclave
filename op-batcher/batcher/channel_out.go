package batcher

import (
	"errors"
	"time"

	"github.com/ethereum-optimism/optimism/op-batcher/batcher"
	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-node/rollup/derive"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum/go-ethereum/core/types"
)

var ErrWithdrawalDetected = errors.New("withdrawal detected")

const minBlockFreshness = 10 * time.Second

func NewChannelOut(cfg batcher.ChannelConfig, rollupCfg *rollup.Config) (derive.ChannelOut, error) {
	co, err := batcher.NewChannelOut(cfg, rollupCfg)
	if err != nil {
		return nil, err
	}
	return &channelOut{
		ChannelOut: co,
	}, nil
}

type channelOut struct {
	derive.ChannelOut
	fullErr            error
	withdrawalDetected bool
}

func (c *channelOut) AddBlock(config *rollup.Config, block *types.Block) (*derive.L1BlockInfo, error) {
	if block.Bloom().Test(predeploys.L2ToL1MessagePasserAddr.Bytes()) {
		c.withdrawalDetected = true
	}
	// If this channel contains a withdrawal, and the block is recent, we can submit the batch immediately.
	// We care about the block freshness because otherwise users could cause the batch submission to slow
	// down significantly by submitting at least 1 withdrawal each block.
	if c.withdrawalDetected && time.Unix(int64(block.Time()), 0).Add(minBlockFreshness).After(time.Now()) {
		c.fullErr = ErrWithdrawalDetected
	}
	return c.ChannelOut.AddBlock(config, block)
}

func (c *channelOut) FullErr() error {
	if c.fullErr != nil {
		return c.fullErr
	}
	return c.ChannelOut.FullErr()
}
