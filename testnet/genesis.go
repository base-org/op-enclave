package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/base-org/op-enclave/bindings"
	altda "github.com/ethereum-optimism/optimism/op-alt-da"
	"github.com/ethereum-optimism/optimism/op-chain-ops/deployer/opcm"
	"github.com/ethereum-optimism/optimism/op-chain-ops/deployer/state"
	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"
	"github.com/ethereum-optimism/optimism/op-chain-ops/genesis"
	"github.com/ethereum-optimism/optimism/op-chain-ops/script"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	"github.com/ethereum-optimism/optimism/op-service/ctxinterrupt"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

var (
	L1URLFlag = &cli.StringFlag{
		Name:     "l1-url",
		Usage:    "URL of an L1 RPC host",
		EnvVars:  []string{"L1_URL"},
		Required: true,
	}
	DeployChainAddressFlag = &cli.StringFlag{
		Name:     "deploy-chain-address",
		Usage:    "Address of the DeployChain contract",
		EnvVars:  []string{"DEPLOY_CHAIN_ADDRESS"},
		Required: true,
	}
	DeployPrivateKeyFlag = &cli.StringFlag{
		Name:     "deploy-private-key",
		Usage:    "Private key of the deployer",
		EnvVars:  []string{"DEPLOY_PRIVATE_KEY"},
		Required: true,
	}
	ConfigPathFlag = &cli.PathFlag{
		Name:    "config-path",
		Usage:   "Path to the config file",
		EnvVars: []string{"CONFIG_PATH"},
	}
	OutputPathFlag = &cli.PathFlag{
		Name:    "output-path",
		Usage:   "Path to the output directory",
		EnvVars: []string{"OUTPUT_PATH"},
		Value:   "./deployments",
	}
	CustomGasTokenFlag = &cli.StringFlag{
		Name:    "gas-token",
		Usage:   "Use a custom gas token",
		EnvVars: []string{"CUSTOM_GAS_TOKEN"},
	}
	DisableProofsFlag = &cli.BoolFlag{
		Name:    "disable-proofs",
		Usage:   "Disable proofs",
		EnvVars: []string{"DISABLE_PROOFS"},
	}
)

var Flags = []cli.Flag{
	L1URLFlag,
	DeployChainAddressFlag,
	DeployPrivateKeyFlag,
	ConfigPathFlag,
	OutputPathFlag,
	CustomGasTokenFlag,
	DisableProofsFlag,
}

type Config struct {
	ProxyAdminOwner            common.Address `json:"proxyAdminOwner"`
	FinalSystemOwner           common.Address `json:"finalSystemOwner"`
	BatchSenderAddress         common.Address `json:"batchSenderAddress"`
	L2OutputOracleProposer     common.Address `json:"l2OutputOracleProposer"`
	P2PSequencerAddress        common.Address `json:"p2pSequencerAddress"`
	BaseFeeVaultRecipient      common.Address `json:"baseFeeVaultRecipient"`
	L1FeeVaultRecipient        common.Address `json:"l1FeeVaultRecipient"`
	SequencerFeeVaultRecipient common.Address `json:"sequencerFeeVaultRecipient"`
}

func (c *Config) Check() error {
	if c.ProxyAdminOwner == (common.Address{}) {
		return errors.New("missing proxy admin owner")
	}
	if c.FinalSystemOwner == (common.Address{}) {
		return errors.New("missing final system owner")
	}
	if c.BatchSenderAddress == (common.Address{}) {
		return errors.New("missing batch sender address")
	}
	if c.L2OutputOracleProposer == (common.Address{}) {
		return errors.New("missing L2 output oracle proposer")
	}
	if c.P2PSequencerAddress == (common.Address{}) {
		return errors.New("missing P2P sequencer address")
	}
	if c.BaseFeeVaultRecipient == (common.Address{}) {
		return errors.New("missing base fee vault recipient")
	}
	if c.L1FeeVaultRecipient == (common.Address{}) {
		return errors.New("missing L1 fee vault recipient")
	}
	if c.SequencerFeeVaultRecipient == (common.Address{}) {
		return errors.New("missing sequencer fee vault recipient")
	}
	return nil

}

func main() {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(Flags)
	app.Name = "genesis"
	app.Usage = "Generate the L2 genesis and deploy the proxies on L1"
	app.Action = Main

	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}

func Main(cliCtx *cli.Context) error {
	l1URL := cliCtx.String(L1URLFlag.Name)
	deployChainContractAddress := common.HexToAddress(cliCtx.String(DeployChainAddressFlag.Name))
	deployPrivateKey := cliCtx.String(DeployPrivateKeyFlag.Name)
	configPath := cliCtx.Path(ConfigPathFlag.Name)
	outputPath := cliCtx.Path(OutputPathFlag.Name)
	customGasToken := cliCtx.String(CustomGasTokenFlag.Name)
	disableProofs := cliCtx.Bool(DisableProofsFlag.Name)

	err := os.MkdirAll(outputPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	deployKey, err := crypto.ToECDSA(common.FromHex(deployPrivateKey))
	if err != nil {
		return fmt.Errorf("failed to parse deploy private key: %w", err)
	}

	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, l1URL)
	if err != nil {
		return fmt.Errorf("failed to connect to the Ethereum client: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to read ChainID: %w", err)
	}

	deployChain, err := bindings.NewDeployChain(deployChainContractAddress, client)
	if err != nil {
		return fmt.Errorf("failed to create DeployChain binding: %w", err)
	}

	l2ChainID, err := rand.Int(rand.Reader, big.NewInt(math.MaxUint32))
	if err != nil {
		return fmt.Errorf("failed to generate L2 chain ID: %w", err)
	}

	l1Addresses, err := deployChain.DeployAddresses(&bind.CallOpts{}, l2ChainID)
	if err != nil {
		return fmt.Errorf("failed to get L2 proxy deploy addresses: %w", err)
	}

	protocolVersions, err := deployChain.ProtocolVersions(&bind.CallOpts{})
	if err != nil {
		return fmt.Errorf("failed to get ProtocolVersions contract address: %w", err)
	}

	batchInboxAddress, err := deployChain.CalculateBatchInbox(&bind.CallOpts{}, l2ChainID)
	if err != nil {
		return fmt.Errorf("failed to calculate BatchInbox address: %w", err)
	}

	if customGasToken != "" {
		log.Info("Checking custom gas token", "address", customGasToken)
		code, err := client.CodeAt(ctx, common.HexToAddress(customGasToken), nil)
		if err != nil {
			return fmt.Errorf("failed to check custom gas token code: %w", err)
		}
		if len(code) == 0 {
			return fmt.Errorf("custom gas token address %s has no code", customGasToken)
		}
	}

	prefix := filepath.Join(outputPath, fmt.Sprintf("%s-%s-", chainID.String(), l2ChainID.String()))
	var genesisConfig Config
	if configPath == "" {
		keysFile := prefix + "keys.json"
		log.Warn("No config file provided, generating new keys", "file", keysFile)

		keys := make(map[string]string)
		names := []string{
			"proxyAdminOwner",
			"finalSystemOwner",
			"batchSender",
			"l2OutputOracleProposer",
			"p2pSequencer",
			"vaultRecipient",
		}
		for _, name := range names {
			key, err := crypto.GenerateKey()
			if err != nil {
				return fmt.Errorf("failed to generate key for %s: %w", name, err)
			}
			keys[name] = common.Bytes2Hex(crypto.FromECDSA(key))
		}
		keysJSON, err := json.MarshalIndent(keys, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal keys: %w", err)
		}
		err = os.WriteFile(keysFile, keysJSON, 0644)
		if err != nil {
			return fmt.Errorf("failed to write keys: %w", err)
		}
		toAddress := func(name string) common.Address {
			return crypto.PubkeyToAddress(crypto.ToECDSAUnsafe(common.FromHex(keys[name])).PublicKey)
		}
		genesisConfig.ProxyAdminOwner = toAddress("proxyAdminOwner")
		genesisConfig.FinalSystemOwner = toAddress("finalSystemOwner")
		genesisConfig.BatchSenderAddress = toAddress("batchSender")
		genesisConfig.L2OutputOracleProposer = toAddress("l2OutputOracleProposer")
		genesisConfig.P2PSequencerAddress = toAddress("p2pSequencer")
		genesisConfig.BaseFeeVaultRecipient = toAddress("vaultRecipient")
		genesisConfig.L1FeeVaultRecipient = toAddress("vaultRecipient")
		genesisConfig.SequencerFeeVaultRecipient = toAddress("vaultRecipient")
	} else {
		configJSON, err := os.ReadFile(configPath)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		err = json.Unmarshal(configJSON, &genesisConfig)
		if err != nil {
			return fmt.Errorf("failed to parse config file %s: %w", configPath, err)
		}
	}
	err = genesisConfig.Check()
	if err != nil {
		return fmt.Errorf("config check failed: %w", err)
	}

	config := state.DefaultDeployConfig()

	// copy config from input
	config.ProxyAdminOwner = genesisConfig.ProxyAdminOwner
	config.FinalSystemOwner = genesisConfig.FinalSystemOwner
	config.P2PSequencerAddress = genesisConfig.P2PSequencerAddress
	config.BatchSenderAddress = genesisConfig.BatchSenderAddress
	config.L2OutputOracleProposer = genesisConfig.L2OutputOracleProposer
	config.BaseFeeVaultRecipient = genesisConfig.BaseFeeVaultRecipient
	config.SequencerFeeVaultRecipient = genesisConfig.SequencerFeeVaultRecipient
	config.L1FeeVaultRecipient = genesisConfig.L1FeeVaultRecipient

	// set up defaults
	config.L1ChainID = chainID.Uint64()
	config.L2ChainID = l2ChainID.Uint64()
	config.EnableGovernance = false
	config.L1BlockTime = 2
	config.L2BlockTime = 1
	config.SequencerWindowSize = 3600 * 6 // to account for L2 block time
	config.BatchInboxAddress = batchInboxAddress
	config.FinalizationPeriodSeconds = 1
	config.UseAltDA = true
	config.DACommitmentType = altda.GenericCommitmentString
	config.DAChallengeWindow = 1
	config.DAResolveWindow = 1
	config.L2OutputOracleChallenger = common.Address{1}
	config.SuperchainConfigGuardian = common.Address{1}
	config.L2OutputOracleSubmissionInterval = 1

	// set up deployed contract addresses
	config.L1StandardBridgeProxy = l1Addresses.L1StandardBridge
	config.L1CrossDomainMessengerProxy = l1Addresses.L1CrossDomainMessenger
	config.L1ERC721BridgeProxy = l1Addresses.L1ERC721Bridge
	config.SystemConfigProxy = l1Addresses.SystemConfig
	config.OptimismPortalProxy = l1Addresses.OptimismPortal
	config.DAChallengeProxy = common.Address{}
	config.ProtocolVersionsProxy = protocolVersions

	// custom token configuration
	if customGasToken != "" {

		config.UseCustomGasToken = true
		config.CustomGasTokenAddress = common.HexToAddress(customGasToken)
	}

	l1Header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to get recent L1 block header: %w", err)
	}

	l1Hash := l1Header.Hash()
	config.L1StartingBlockTag = &genesis.MarshalableRPCBlockNumberOrHash{
		BlockHash: &l1Hash,
	}

	err = config.Check(log.Root())
	if err != nil {
		return fmt.Errorf("deploy config check failed: %w", err)
	}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	err = os.WriteFile(prefix+"deploy-config.json", configJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	l2Genesis, err := genesis.NewL2Genesis(&config, l1Header)
	if err != nil {
		return fmt.Errorf("failed to create L2 genesis: %w", err)
	}

	log.Info("Building Optimism contracts")
	cmd := exec.Command("forge", "build")
	cmd.Dir = "./lib/optimism/packages/contracts-bedrock"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to build contracts: %w", err)
	}

	foundryArtifacts := foundry.OpenArtifactsDir("./lib/optimism/packages/contracts-bedrock/forge-artifacts")
	sourceMap := foundry.NewSourceMapFS(os.DirFS("./lib/optimism/packages/contracts-bedrock"))
	l2Host := createL2(log.Root(), foundryArtifacts, sourceMap, &config, l2Genesis.Timestamp)
	if err := l2Host.EnableCheats(); err != nil {
		return fmt.Errorf("failed to enable cheats: %w", err)
	}
	if err := opcm.L2Genesis(l2Host, &opcm.L2GenesisInput{
		L1Deployments: opcm.L1Deployments{
			L1CrossDomainMessengerProxy: config.L1CrossDomainMessengerProxy,
			L1StandardBridgeProxy:       config.L1StandardBridgeProxy,
			L1ERC721BridgeProxy:         config.L1ERC721BridgeProxy,
		},
		L2Config: config.L2InitializationConfig,
	}); err != nil {
		return fmt.Errorf("failed L2 genesis: %w", err)
	}
	allocs, err := l2Host.StateDump()
	if err != nil {
		return fmt.Errorf("failed to dump L1 state: %w", err)
	}
	if err := ensureNoDeployed(allocs, sysGenesisDeployer); err != nil {
		return fmt.Errorf("unexpected deployed account content by L2 genesis deployer: %w", err)
	}

	l2Genesis.Alloc = allocs.Accounts
	addPredeploys(l2Genesis.Alloc)
	genesisBlock := l2Genesis.ToBlock()

	genesisJSON, err := json.MarshalIndent(l2Genesis, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal genesis: %w", err)
	}
	err = os.WriteFile(prefix+"genesis.json", genesisJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write genesis: %w", err)
	}

	rollupConfig, err := config.RollupConfig(l1Header, genesisBlock.Hash(), 0)
	if err != nil {
		return fmt.Errorf("failed to create rollup config: %w", err)
	}
	rollupConfigJSON, err := json.MarshalIndent(rollupConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal rollup config: %w", err)
	}
	err = os.WriteFile(prefix+"rollup-config.json", rollupConfigJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write rollup config: %w", err)
	}

	signer := types.LatestSignerForChainID(chainID)
	opts := &bind.TransactOpts{
		From: crypto.PubkeyToAddress(deployKey.PublicKey),
		Signer: func(_ common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return types.SignTx(tx, signer, deployKey)
		},
	}

	genesisCfg := bindings.DeployChainGenesisConfiguration{
		L1Number:    l1Header.Number.Uint64(),
		L2Hash:      genesisBlock.Hash(),
		L2StateRoot: genesisBlock.Root(),
		L2Time:      l2Genesis.Timestamp,
	}
	gasConfig := bindings.DeployChainGasConfiguration{
		BasefeeScalar:     config.GasPriceOracleBaseFeeScalar,
		BlobbasefeeScalar: config.GasPriceOracleBlobBaseFeeScalar,
		GasLimit:          uint64(config.L2GenesisBlockGasLimit),
		GasToken:          config.CustomGasTokenAddress,
	}
	addressConfig := bindings.DeployChainAddressConfiguration{
		Batcher:           config.BatchSenderAddress,
		Proposer:          config.L2OutputOracleProposer,
		UnsafeBlockSigner: config.P2PSequencerAddress,
	}

	tx, err := deployChain.Deploy(
		opts,
		l2ChainID,
		genesisCfg,
		gasConfig,
		addressConfig,
		!disableProofs,
	)
	if err != nil {
		return fmt.Errorf("failed to deploy proxies: %w", err)
	}

	log.Info("Deployed proxies", "tx", tx.Hash().Hex())

	receipt, err := waitForConfirmation(ctx, client, tx.Hash())
	if err != nil {
		return fmt.Errorf("error waiting for confirmation: %w", err)
	}
	var deploy *bindings.DeployChainDeploy
	for _, l := range receipt.Logs {
		deploy, err = deployChain.ParseDeploy(*l)
		if err == nil {
			break
		}
	}
	if deploy == nil {
		return errors.New("failed to parse Deploy event")
	}

	deployJSON, err := json.MarshalIndent(deploy, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal deployed info: %w", err)
	}
	err = os.WriteFile(prefix+"deployed.json", deployJSON, 0644)
	if err != nil {
		return fmt.Errorf("failed to write deployed info: %w", err)
	}

	return nil
}

func waitForConfirmation(ctx context.Context, client *ethclient.Client, tx common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(ctx, tx)
		if errors.Is(err, ethereum.NotFound) {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(1 * time.Second):
			}
		} else if err != nil {
			return nil, err
		} else if receipt.Status != types.ReceiptStatusSuccessful {
			return nil, errors.New("unsuccessful receipt status")
		} else {
			return receipt, nil
		}
	}
}

var sysGenesisDeployer = common.Address(crypto.Keccak256([]byte("System genesis deployer"))[12:])

func createL2(logger log.Logger, fa *foundry.ArtifactsFS, srcFS *foundry.SourceMapFS, deployConfig *genesis.DeployConfig, genesisTimestamp uint64) *script.Host {
	l2Context := script.Context{
		ChainID:      new(big.Int).SetUint64(deployConfig.L2ChainID),
		Sender:       sysGenesisDeployer,
		Origin:       sysGenesisDeployer,
		FeeRecipient: common.Address{},
		GasLimit:     script.DefaultFoundryGasLimit,
		BlockNum:     uint64(deployConfig.L2GenesisBlockNumber),
		Timestamp:    genesisTimestamp,
		PrevRandao:   deployConfig.L2GenesisBlockMixHash,
		BlobHashes:   nil,
	}
	l2Host := script.NewHost(logger.New("role", "l2", "chain", deployConfig.L2ChainID), fa, srcFS, l2Context)
	l2Host.SetEnvVar("OUTPUT_MODE", "none") // we don't use the cheatcode, but capture the state outside of EVM execution
	l2Host.SetEnvVar("FORK", "granite")     // latest fork
	return l2Host
}

func ensureNoDeployed(allocs *foundry.ForgeAllocs, deployer common.Address) error {
	// Sanity check we have no deploy output that's not meant to be there.
	for i := uint64(0); i <= allocs.Accounts[deployer].Nonce; i++ {
		addr := crypto.CreateAddress(deployer, i)
		if _, ok := allocs.Accounts[addr]; ok {
			return fmt.Errorf("system deployer output %s (deployed with nonce %d) was not cleaned up", addr, i)
		}
	}
	// Don't include the deployer account
	delete(allocs.Accounts, deployer)
	return nil
}
