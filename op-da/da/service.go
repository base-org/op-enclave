package da

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync/atomic"

	"github.com/base/op-enclave/op-da/flags"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

func Main(cliCtx *cli.Context, _ context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	cfg := NewConfig(cliCtx)

	l := oplog.NewLogger(oplog.AppOut(cliCtx), cfg.LogConfig)
	oplog.SetGlobalLogHandler(l.Handler())
	opservice.ValidateEnvVars(flags.EnvVarPrefix, flags.Flags, l)

	l.Info("Initializing alt-DA server")
	return ServiceFromCLIConfig(cliCtx.Context, cfg, l)
}

func ServiceFromCLIConfig(ctx context.Context, cfg *CLIConfig, l log.Logger) (cliapp.Lifecycle, error) {
	store, err := newStore(cfg)
	if err != nil {
		return nil, err
	}

	server := altda.NewDAServer("", cfg.Port, store, l, true)

	return &service{
		server: server,
	}, nil
}

type service struct {
	server  *altda.DAServer
	stopped atomic.Bool
}

func (s *service) Start(ctx context.Context) error {
	return s.server.Start()
}

func (s *service) Stop(ctx context.Context) error {
	if s.stopped.Swap(true) {
		return nil
	}
	return s.server.Stop()
}

func (s *service) Stopped() bool {
	return s.stopped.Load()
}

func newStore(cfg *CLIConfig) (altda.KVStore, error) {
	u, err := url.Parse(cfg.DAURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DA URL: %w", err)
	}
	switch u.Scheme {
	case "s3":
		return NewS3store(u.Host, u.Path), nil
	case "file":
		err = os.MkdirAll(u.Path, os.FileMode(0755))
		if err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
		return NewFilestore(u.Path), nil
	default:
		return nil, fmt.Errorf("unsupported DA scheme: %s", u.Scheme)
	}
}
