// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import {Deploy} from "@eth-optimism-bedrock/scripts/deploy/Deploy.s.sol";
import {DeployConfig} from "@eth-optimism-bedrock/scripts/deploy/DeployConfig.s.sol";
import {Types} from "@eth-optimism-bedrock/scripts/libraries/Types.sol";
import {ChainAssertions} from "@eth-optimism-bedrock/scripts/deploy/ChainAssertions.sol";
import {SystemConfig} from "@eth-optimism-bedrock/src/L1/SystemConfig.sol";
import {ISystemConfig} from "@eth-optimism-bedrock/src/L1/interfaces/ISystemConfig.sol";
import {ISuperchainConfig} from "@eth-optimism-bedrock/src/L1/interfaces/ISuperchainConfig.sol";
import {SuperchainConfig} from "@eth-optimism-bedrock/src/L1/SuperchainConfig.sol";
import {OptimismPortal} from "@eth-optimism-bedrock/src/L1/OptimismPortal.sol";
import {IL2OutputOracle} from "@eth-optimism-bedrock/src/L1/interfaces/IL2OutputOracle.sol";
import {ProtocolVersions} from "@eth-optimism-bedrock/src/L1/ProtocolVersions.sol";
import {Portal} from "src/Portal.sol";
import {OutputOracle} from "src/OutputOracle.sol";
import {OwnerConfig} from "src/OwnerConfig.sol";
import {SystemConfigOwnable} from "src/SystemConfigOwnable.sol";
import {SystemConfigGlobal} from "src/SystemConfigGlobal.sol";
import {DeployChain} from "src/DeployChain.sol";
import {Constants} from "@eth-optimism-bedrock/src/libraries/Constants.sol";
import {ResourceMetering} from "@eth-optimism-bedrock/src/L1/ResourceMetering.sol";
import {IResourceMetering} from "@eth-optimism-bedrock/src/L1/interfaces/IResourceMetering.sol";
import {ICertManager} from "@nitro-validator/ICertManager.sol";

import {console2 as console} from "forge-std/console2.sol";

contract DeploySystem is Deploy {
    function deploy() public {
        console.log("start of L1 Deploy!");
        deploySafe("SystemOwnerSafe");
        console.log("deployed Safe!");

        setupAdmin();

        setupSuperchain();
        console.log("set up superchain!");

        setupSystemConfigGlobal();
        console.log("set up system config global!");

        setupOpChain2();
        console.log("set up op chain!");

        setupChainDeploy();
        console.log("set chain deploy!");
    }

    function setupSystemConfigGlobal() public {
        console.log("Setting up SystemConfigGlobal");

        checkCertManager();

        deployERC1967Proxy("SystemConfigGlobalProxy");
        deploySystemConfigGlobal();
        initializeSystemConfigGlobal();
    }

    function setupOpChain2() public {
        console.log("Deploying OP Chain");

        // Ensure that the requisite contracts are deployed
        mustGetAddress("SuperchainConfigProxy");
        mustGetAddress("SystemConfigGlobalProxy");
        mustGetAddress("SystemOwnerSafe");
        mustGetAddress("AddressManager");
        mustGetAddress("ProxyAdmin");

        deployProxies2();
        deployImplementations2();
        initializeImplementations2();
    }

    function setupChainDeploy() public {
        console.log("Setting up Chain Deploy");

        deployDeployChain();
    }

    function deployProxies2() public {
        console.log("Deploying proxies");

        deployERC1967Proxy("OptimismPortalProxy");
        deployERC1967Proxy("SystemConfigProxy");
        deployL1StandardBridgeProxy();
        deployL1CrossDomainMessengerProxy();
        deployERC1967Proxy("OptimismMintableERC20FactoryProxy");
        deployERC1967Proxy("L1ERC721BridgeProxy");
        deployERC1967Proxy("L2OutputOracleProxy");

        transferAddressManagerOwnership(); // to the ProxyAdmin

        // fake these, required for calls to `_proxies()`
        save("DisputeGameFactoryProxy", address(1));
        save("DelayedWETHProxy", address(1));
        save("PermissionedDelayedWETHProxy", address(1));
        save("AnchorStateRegistryProxy", address(1));
    }

    function deployImplementations2() public {
        console.log("Deploying implementations");
        deployL1CrossDomainMessenger();
        deployOptimismMintableERC20Factory();
        deploySystemConfigOwnable();
        deployL1StandardBridge();
        deployL1ERC721Bridge();
        deployPortal();
        deployOutputOracle();
    }

    function initializeImplementations2() public {
        console.log("Initializing implementations");
        initializePortal();
        initializeSystemConfigOwnable();
        initializeL1StandardBridge();
        initializeL1ERC721Bridge();
        initializeOptimismMintableERC20Factory();
        initializeL1CrossDomainMessenger();
        initializeOutputOracle();
    }

    function checkCertManager() public {
        console.log("Retrieving CertManager deploy");
        string memory deploymentOutfile =
            string.concat(vm.projectRoot(), "/deployments/", vm.toString(block.chainid), "-certmanager.json");
        address certManagerAddress = vm.parseJsonAddress(vm.readFile(deploymentOutfile), ".CertManager");
        save("CertManager", certManagerAddress);

        bytes memory parent =
            hex"3082021130820196a003020102021100f93175681b90afe11d46ccb4e4e7f856300a06082a8648ce3d0403033049310b3009060355040613025553310f300d060355040a0c06416d617a6f6e310c300a060355040b0c03415753311b301906035504030c126177732e6e6974726f2d656e636c61766573301e170d3139313032383133323830355a170d3439313032383134323830355a3049310b3009060355040613025553310f300d060355040a0c06416d617a6f6e310c300a060355040b0c03415753311b301906035504030c126177732e6e6974726f2d656e636c617665733076301006072a8648ce3d020106052b8104002203620004fc0254eba608c1f36870e29ada90be46383292736e894bfff672d989444b5051e534a4b1f6dbe3c0bc581a32b7b176070ede12d69a3fea211b66e752cf7dd1dd095f6f1370f4170843d9dc100121e4cf63012809664487c9796284304dc53ff4a3423040300f0603551d130101ff040530030101ff301d0603551d0e041604149025b50dd90547e796c396fa729dcf99a9df4b96300e0603551d0f0101ff040403020186300a06082a8648ce3d0403030369003066023100a37f2f91a1c9bd5ee7b8627c1698d255038e1f0343f95b63a9628c3d39809545a11ebcbf2e3b55d8aeee71b4c3d6adf3023100a2f39b1605b27028a5dd4ba069b5016e65b4fbde8fe0061d6a53197f9cdaf5d943bc61fc2beb03cb6fee8d2302f3dff6";
        bytes memory cert =
            hex"308202bf30820244a00302010202100b93e39c65609c59e8144a2ad34ba3a0300a06082a8648ce3d0403033049310b3009060355040613025553310f300d060355040a0c06416d617a6f6e310c300a060355040b0c03415753311b301906035504030c126177732e6e6974726f2d656e636c61766573301e170d3234313132333036333235355a170d3234313231333037333235355a3064310b3009060355040613025553310f300d060355040a0c06416d617a6f6e310c300a060355040b0c034157533136303406035504030c2d353133623665666332313639303264372e75732d656173742d312e6177732e6e6974726f2d656e636c617665733076301006072a8648ce3d020106052b8104002203620004ee78108039725a03e0b63a5d7d1244f6294eb7631f305e360997c8e5c06c779f23cfaeb64cb9aeac8a031bfac9f4dafc3621b4367f003c08c0ce410c2118396cc5d56ec4e92e1b17f9709b2bffcef462f7bcb97d6ca11325c4a30156c9720de7a381d53081d230120603551d130101ff040830060101ff020102301f0603551d230418301680149025b50dd90547e796c396fa729dcf99a9df4b96301d0603551d0e041604142b3d75d274a3cdd61b2c13f539e08c960ce757dd300e0603551d0f0101ff040403020186306c0603551d1f046530633061a05fa05d865b687474703a2f2f6177732d6e6974726f2d656e636c617665732d63726c2e73332e616d617a6f6e6177732e636f6d2f63726c2f61623439363063632d376436332d343262642d396539662d3539333338636236376638342e63726c300a06082a8648ce3d0403030369003066023100fce7a6c2b38e0a8ebf0d28348d74463458b84bfe8b2b95315dd4da665e8e83d4ab911852a4e92a8263ecf571d2df3b89023100ab92be511136be76aa313018f9f4825eaad602d0342d268e6da632767f68f55f761fa9fd2a7ee716c481c67f26e3f8f4";

        uint256 timestamp = vm.getBlockTimestamp();
        vm.warp(1732580000);
        ICertManager(certManagerAddress).verifyCert(cert, false, keccak256(parent));
        vm.warp(timestamp);
    }

    function deploySystemConfigOwnable() public broadcast returns (address addr_) {
        console.log("Deploying OwnerConfig");
        OwnerConfig ownerConfig = new OwnerConfig{salt: _implSalt()}(cfg.finalSystemOwner());
        save("OwnerConfig", address(ownerConfig));

        console.log("Deploying SystemConfig implementation");
        addr_ = address(new SystemConfigOwnable{salt: _implSalt()}(ownerConfig));
        save("SystemConfig", addr_);
        console.log("SystemConfig deployed at %s", addr_);

        // Override the `SystemConfig` contract to the deployed implementation. This is necessary
        // to check the `SystemConfig` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.SystemConfig = addr_;
        checkSystemConfig({_contracts: contracts, _cfg: cfg, _isProxy: false});
    }

    function deploySystemConfigGlobal() public broadcast returns (address addr_) {
        console.log("Deploying SystemConfigGlobal implementation");
        addr_ = address(new SystemConfigGlobal{salt: _implSalt()}(ICertManager(mustGetAddress("CertManager"))));
        save("SystemConfigGlobal", addr_);
        console.log("SystemConfigGlobal deployed at %s", addr_);
    }

    function deployPortal() public broadcast returns (address addr_) {
        console.log("Deploying OptimismPortal implementation");
        addr_ = address(new Portal{salt: _implSalt()}());
        save("OptimismPortal", addr_);
        console.log("OptimismPortal deployed at %s", addr_);

        // Override the `OptimismPortal` contract to the deployed implementation. This is necessary
        // to check the `OptimismPortal` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.OptimismPortal = addr_;
        ChainAssertions.checkOptimismPortal({_contracts: contracts, _cfg: cfg, _isProxy: false});
    }

    function deployOutputOracle() public broadcast returns (address addr_) {
        SystemConfigGlobal systemConfigGlobal = SystemConfigGlobal(mustGetAddress("SystemConfigGlobalProxy"));

        console.log("Deploying L2OutputOracle implementation");
        OutputOracle oracle = new OutputOracle{salt: _implSalt()}(
            systemConfigGlobal,
            1000 // TODO move to config
        );

        save("L2OutputOracle", address(oracle));
        console.log("L2OutputOracle deployed at %s", address(oracle));

        // Override the `L2OutputOracle` contract to the deployed implementation. This is necessary
        // to check the `L2OutputOracle` implementation alongside dependent contracts, which
        // are always proxies.
        Types.ContractSet memory contracts = _proxiesUnstrict();
        contracts.L2OutputOracle = address(oracle);
        checkL2OutputOracle({_contracts: contracts, _cfg: cfg, _isProxy: false});

        addr_ = address(oracle);
    }

    function deployDeployChain() public broadcast returns (address addr_) {
        console.log("Deploying DeployChain implementation");
        DeployChain deployChain = new DeployChain{salt: _implSalt()}({
            _proxyAdmin: mustGetAddress("ProxyAdmin"),
            _optimismPortal: mustGetAddress("OptimismPortalProxy"),
            _systemConfig: mustGetAddress("SystemConfigProxy"),
            _l1StandardBridge: mustGetAddress("L1StandardBridgeProxy"),
            _l1ERC721Bridge: mustGetAddress("L1ERC721BridgeProxy"),
            _optimismMintableERC20Factory: mustGetAddress("OptimismMintableERC20FactoryProxy"),
            _l1CrossDomainMessenger: mustGetAddress("L1CrossDomainMessengerProxy"),
            _l2OutputOracle: mustGetAddress("L2OutputOracleProxy"),
            _superchainConfig: mustGetAddress("SuperchainConfigProxy"),
            _protocolVersions: mustGetAddress("ProtocolVersionsProxy")
        });

        save("DeployChain", address(deployChain));
        console.log("DeployChain deployed at %s", address(deployChain));

        addr_ = address(deployChain);
    }

    function initializeSystemConfigOwnable() public broadcast {
        console.log("Upgrading and initializing SystemConfig proxy");
        address systemConfigProxy = mustGetAddress("SystemConfigProxy");
        address systemConfig = mustGetAddress("SystemConfig");

        bytes32 batcherHash = bytes32(uint256(uint160(cfg.batchSenderAddress())));

        _upgradeAndCallViaSafe({
            _proxy: payable(systemConfigProxy),
            _implementation: systemConfig,
            _innerCallData: abi.encodeCall(
                SystemConfigOwnable.initialize,
                (
                    cfg.basefeeScalar(),
                    cfg.blobbasefeeScalar(),
                    batcherHash,
                    uint64(cfg.l2GenesisBlockGasLimit()),
                    cfg.p2pSequencerAddress(),
                    Constants.DEFAULT_RESOURCE_CONFIG(),
                    cfg.batchInboxAddress(),
                    address(0),
                    SystemConfig.Addresses({
                        l1CrossDomainMessenger: mustGetAddress("L1CrossDomainMessengerProxy"),
                        l1ERC721Bridge: mustGetAddress("L1ERC721BridgeProxy"),
                        l1StandardBridge: mustGetAddress("L1StandardBridgeProxy"),
                        disputeGameFactory: mustGetAddress("DisputeGameFactoryProxy"),
                        optimismPortal: mustGetAddress("OptimismPortalProxy"),
                        optimismMintableERC20Factory: mustGetAddress("OptimismMintableERC20FactoryProxy"),
                        gasPayingToken: Constants.ETHER
                    })
                )
            )
        });

        SystemConfig config = SystemConfig(systemConfigProxy);
        string memory version = config.version();
        console.log("SystemConfig version: %s", version);

        checkSystemConfig({_contracts: _proxies(), _cfg: cfg, _isProxy: true});
    }

    function initializeSystemConfigGlobal() public broadcast {
        console.log("Upgrading and initializing SystemConfigGlobal proxy");
        address systemConfigGlobalProxy = mustGetAddress("SystemConfigGlobalProxy");
        address systemConfigGlobal = mustGetAddress("SystemConfigGlobal");

        _upgradeAndCallViaSafe({
            _proxy: payable(systemConfigGlobalProxy),
            _implementation: systemConfigGlobal,
            _innerCallData: abi.encodeWithSelector(SystemConfigGlobal.initialize.selector, cfg.finalSystemOwner())
        });

        SystemConfigGlobal config = SystemConfigGlobal(systemConfigGlobalProxy);
        string memory version = config.version();
        console.log("SystemConfigGlobal version: %s", version);
    }

    function initializePortal() public broadcast {
        console.log("Upgrading and initializing OptimismPortal proxy");
        address optimismPortalProxy = mustGetAddress("OptimismPortalProxy");
        address optimismPortal = mustGetAddress("OptimismPortal");
        address l2OutputOracleProxy = mustGetAddress("L2OutputOracleProxy");
        address systemConfigProxy = mustGetAddress("SystemConfigProxy");
        address superchainConfigProxy = mustGetAddress("SuperchainConfigProxy");

        _upgradeAndCallViaSafe({
            _proxy: payable(optimismPortalProxy),
            _implementation: optimismPortal,
            _innerCallData: abi.encodeCall(
                OptimismPortal.initialize,
                (
                    IL2OutputOracle(l2OutputOracleProxy),
                    ISystemConfig(systemConfigProxy),
                    ISuperchainConfig(superchainConfigProxy)
                )
            )
        });

        Portal portal = Portal(payable(optimismPortalProxy));
        string memory version = portal.version();
        console.log("OptimismPortal version: %s", version);

        ChainAssertions.checkOptimismPortal({_contracts: _proxies(), _cfg: cfg, _isProxy: true});
    }

    function initializeOutputOracle() public broadcast {
        console.log("Upgrading and initializing OutputOracle proxy");
        address l2OutputOracleProxy = mustGetAddress("L2OutputOracleProxy");
        address l2OutputOracle = mustGetAddress("L2OutputOracle");
        address systemConfigProxy = mustGetAddress("SystemConfigProxy");

        _upgradeAndCallViaSafe({
            _proxy: payable(l2OutputOracleProxy),
            _implementation: l2OutputOracle,
            _innerCallData: abi.encodeCall(OutputOracle.initialize, (SystemConfigOwnable(systemConfigProxy), 0, 0, false))
        });

        OutputOracle oracle = OutputOracle(l2OutputOracleProxy);
        string memory version = oracle.version();
        console.log("L2OutputOracle version: %s", version);

        checkL2OutputOracle({_contracts: _proxies(), _cfg: cfg, _isProxy: true});
    }

    function checkL2OutputOracle(Types.ContractSet memory _contracts, DeployConfig _cfg, bool _isProxy) internal view {
        console.log("Running chain assertions on the L2OutputOracle");
        OutputOracle oracle = OutputOracle(_contracts.L2OutputOracle);

        // Check that the contract is initialized
        ChainAssertions.assertSlotValueIsOne({_contractAddress: address(oracle), _slot: 0, _offset: 0});

        if (_isProxy) {
            require(oracle.proposer() == _cfg.l2OutputOracleProposer());
        } else {
            require(oracle.proposer() == address(0));
        }
    }

    function checkSystemConfig(Types.ContractSet memory _contracts, DeployConfig _cfg, bool _isProxy) internal view {
        console.log("Running chain assertions on the SystemConfig");
        SystemConfig config = SystemConfig(_contracts.SystemConfig);

        // Check that the contract is initialized
        ChainAssertions.assertSlotValueIsOne({_contractAddress: address(config), _slot: 0, _offset: 0});

        IResourceMetering.ResourceConfig memory resourceConfig = config.resourceConfig();

        if (_isProxy) {
            require(config.owner() == _cfg.finalSystemOwner());
            require(config.basefeeScalar() == _cfg.basefeeScalar());
            require(config.blobbasefeeScalar() == _cfg.blobbasefeeScalar());
            require(config.batcherHash() == bytes32(uint256(uint160(_cfg.batchSenderAddress()))));
            require(config.gasLimit() == uint64(_cfg.l2GenesisBlockGasLimit()));
            require(config.unsafeBlockSigner() == _cfg.p2pSequencerAddress());
            require(config.scalar() >> 248 == 1);
            // Check _config
            IResourceMetering.ResourceConfig memory rconfig = Constants.DEFAULT_RESOURCE_CONFIG();
            require(resourceConfig.maxResourceLimit == rconfig.maxResourceLimit);
            require(resourceConfig.elasticityMultiplier == rconfig.elasticityMultiplier);
            require(resourceConfig.baseFeeMaxChangeDenominator == rconfig.baseFeeMaxChangeDenominator);
            require(resourceConfig.systemTxMaxGas == rconfig.systemTxMaxGas);
            require(resourceConfig.minimumBaseFee == rconfig.minimumBaseFee);
            require(resourceConfig.maximumBaseFee == rconfig.maximumBaseFee);
            // Depends on start block being set to 0 in `initialize`
            uint256 cfgStartBlock = _cfg.systemConfigStartBlock();
            require(config.startBlock() == (cfgStartBlock == 0 ? block.number : cfgStartBlock));
            require(config.batchInbox() == _cfg.batchInboxAddress());
            // Check _addresses
            require(config.l1CrossDomainMessenger() == _contracts.L1CrossDomainMessenger);
            require(config.l1ERC721Bridge() == _contracts.L1ERC721Bridge);
            require(config.l1StandardBridge() == _contracts.L1StandardBridge);
            require(config.disputeGameFactory() == _contracts.DisputeGameFactory);
            require(config.optimismPortal() == _contracts.OptimismPortal);
            require(config.optimismMintableERC20Factory() == _contracts.OptimismMintableERC20Factory);
        } else {
            require(config.owner() == _cfg.finalSystemOwner());
            require(config.overhead() == 0);
            require(config.scalar() == uint256(0x01) << 248); // version 1
            require(config.basefeeScalar() == 0);
            require(config.blobbasefeeScalar() == 0);
            require(config.batcherHash() == bytes32(0));
            require(config.gasLimit() == 1);
            require(config.unsafeBlockSigner() == address(0));
            // Check _config
            require(resourceConfig.maxResourceLimit == 1);
            require(resourceConfig.elasticityMultiplier == 1);
            require(resourceConfig.baseFeeMaxChangeDenominator == 2);
            require(resourceConfig.systemTxMaxGas == 0);
            require(resourceConfig.minimumBaseFee == 0);
            require(resourceConfig.maximumBaseFee == 0);
            // Check _addresses
            require(config.startBlock() == type(uint256).max);
            require(config.batchInbox() == address(0));
            require(config.l1CrossDomainMessenger() == address(0));
            require(config.l1ERC721Bridge() == address(0));
            require(config.l1StandardBridge() == address(0));
            require(config.disputeGameFactory() == address(0));
            require(config.optimismPortal() == address(0));
            require(config.optimismMintableERC20Factory() == address(0));
        }
    }
}
