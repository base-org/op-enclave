package main

import (
	"context"
	"os"

	"github.com/base/op-enclave/op-da/da"
	"github.com/base/op-enclave/op-da/flags"
	"github.com/ethereum-optimism/optimism/op-service/ctxinterrupt"

	"github.com/urfave/cli/v2"

	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum/log"
)

func main() {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(flags.Flags)
	app.Name = "op-da"
	app.Usage = "Alt DA server"
	app.Action = cliapp.LifecycleCmd(da.Main)

	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}
