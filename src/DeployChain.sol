// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {ResolvingProxyFactory} from "./ResolvingProxyFactory.sol";
import {Portal} from "./Portal.sol";
import {OutputOracle} from "./OutputOracle.sol";
import {SystemConfigOwnable} from "./SystemConfigOwnable.sol";
import {SystemConfig} from "@eth-optimism-bedrock/src/L1/SystemConfig.sol";
import {ISystemConfig} from "@eth-optimism-bedrock/src/L1/interfaces/ISystemConfig.sol";
import {OptimismPortal} from "@eth-optimism-bedrock/src/L1/OptimismPortal.sol";
import {IOptimismPortal} from "@eth-optimism-bedrock/src/L1/interfaces/IOptimismPortal.sol";
import {ISuperchainConfig} from "@eth-optimism-bedrock/src/L1/interfaces/ISuperchainConfig.sol";
import {L2OutputOracle} from "@eth-optimism-bedrock/src/L1/L2OutputOracle.sol";
import {L1StandardBridge} from "@eth-optimism-bedrock/src/L1/L1StandardBridge.sol";
import {L1CrossDomainMessenger} from "@eth-optimism-bedrock/src/L1/L1CrossDomainMessenger.sol";
import {ICrossDomainMessenger} from "@eth-optimism-bedrock/src/universal/interfaces/ICrossDomainMessenger.sol";
import {L1ERC721Bridge} from "@eth-optimism-bedrock/src/L1/L1ERC721Bridge.sol";
import {OptimismMintableERC20Factory} from "@eth-optimism-bedrock/src/universal/OptimismMintableERC20Factory.sol";
import {ResourceMetering} from "@eth-optimism-bedrock/src/L1/ResourceMetering.sol";
import {Hashing} from "@eth-optimism-bedrock/src/libraries/Hashing.sol";
import {Types} from "@eth-optimism-bedrock/src/libraries/Types.sol";
import {Constants} from "@eth-optimism-bedrock/src/libraries/Constants.sol";

contract DeployChain {
    struct DeployAddresses {
        address l2OutputOracle;
        address systemConfig;
        address optimismPortal;
        address l1CrossDomainMessenger;
        address l1StandardBridge;
        address l1ERC721Bridge;
        address optimismMintableERC20Factory;
    }

    struct GenesisConfiguration {
        uint64 l1Number;
        bytes32 l2Hash;
        bytes32 l2StateRoot;
        uint64 l2Time;
    }

    struct GasConfiguration {
        uint32 basefeeScalar;
        uint32 blobbasefeeScalar;
        uint64 gasLimit;
        address gasToken;
    }

    struct AddressConfiguration {
        address batcher;
        address proposer;
        address unsafeBlockSigner;
    }

    struct Hashes {
        bytes32 configHash;
        bytes32 genesisOutputRoot;
    }

    event Deploy(
        uint256 indexed chainID, bytes32 configHash, bytes32 outputRoot, address batchInbox, DeployAddresses addresses
    );

    bytes32 public constant MESSAGE_PASSER_STORAGE_HASH =
        0x8ed4baae3a927be3dea54996b4d5899f8c01e7594bf50b17dc1e741388ce3d12;

    address public immutable proxyAdmin;
    address public immutable optimismPortal;
    address public immutable systemConfig;
    address public immutable l1StandardBridge;
    address public immutable l1ERC721Bridge;
    address public immutable optimismMintableERC20Factory;
    address public immutable l1CrossDomainMessenger;
    address public immutable l2OutputOracle;
    address public immutable superchainConfig;
    address public immutable protocolVersions;

    constructor(
        address _proxyAdmin,
        address _optimismPortal,
        address _systemConfig,
        address _l1StandardBridge,
        address _l1ERC721Bridge,
        address _optimismMintableERC20Factory,
        address _l1CrossDomainMessenger,
        address _l2OutputOracle,
        address _superchainConfig,
        address _protocolVersions
    ) {
        proxyAdmin = _proxyAdmin;
        optimismPortal = _optimismPortal;
        systemConfig = _systemConfig;
        l1StandardBridge = _l1StandardBridge;
        l1ERC721Bridge = _l1ERC721Bridge;
        optimismMintableERC20Factory = _optimismMintableERC20Factory;
        l1CrossDomainMessenger = _l1CrossDomainMessenger;
        l2OutputOracle = _l2OutputOracle;
        superchainConfig = _superchainConfig;
        protocolVersions = _protocolVersions;
    }

    function deployAddresses(uint256 chainID) external view returns (DeployAddresses memory) {
        bytes32 salt = keccak256(abi.encodePacked(chainID));
        return DeployAddresses({
            l2OutputOracle: ResolvingProxyFactory.proxyAddress(l2OutputOracle, proxyAdmin, salt),
            systemConfig: ResolvingProxyFactory.proxyAddress(systemConfig, proxyAdmin, salt),
            optimismPortal: ResolvingProxyFactory.proxyAddress(optimismPortal, proxyAdmin, salt),
            l1CrossDomainMessenger: ResolvingProxyFactory.proxyAddress(l1CrossDomainMessenger, proxyAdmin, salt),
            l1StandardBridge: ResolvingProxyFactory.proxyAddress(l1StandardBridge, proxyAdmin, salt),
            l1ERC721Bridge: ResolvingProxyFactory.proxyAddress(l1ERC721Bridge, proxyAdmin, salt),
            optimismMintableERC20Factory: ResolvingProxyFactory.proxyAddress(optimismMintableERC20Factory, proxyAdmin, salt)
        });
    }

    function deploy(
        uint256 chainID,
        GenesisConfiguration memory genesisConfig,
        GasConfiguration memory gasConfig,
        AddressConfiguration memory addressConfig
    ) external {
        DeployAddresses memory addresses = setupProxies(chainID);

        Hashes memory hashes = calculateHashes(chainID, genesisConfig, gasConfig, addressConfig.batcher, addresses);

        address batchInbox = calculateBatchInbox(chainID);

        initializeProxies(gasConfig, addressConfig, batchInbox, hashes, addresses);

        emit Deploy({
            chainID: chainID,
            configHash: hashes.configHash,
            outputRoot: hashes.genesisOutputRoot,
            batchInbox: batchInbox,
            addresses: addresses
        });
    }

    function calculateBatchInbox(uint256 chainID) public pure returns (address) {
        uint256 reverse = 0;
        for (; chainID > 0; chainID /= 10) {
            reverse = (reverse * 10) + (chainID % 10);
        }
        uint256 base16 = 0;
        for (; reverse > 0; reverse /= 10) {
            base16 = (base16 << 4) | (reverse % 10);
        }
        return address(uint160(base16 | (0xff << 152)));
    }

    function setupProxies(uint256 chainID) internal returns (DeployAddresses memory) {
        bytes32 salt = keccak256(abi.encodePacked(chainID));
        return DeployAddresses({
            l2OutputOracle: ResolvingProxyFactory.setupProxy(l2OutputOracle, proxyAdmin, salt),
            systemConfig: ResolvingProxyFactory.setupProxy(systemConfig, proxyAdmin, salt),
            optimismPortal: ResolvingProxyFactory.setupProxy(optimismPortal, proxyAdmin, salt),
            l1CrossDomainMessenger: ResolvingProxyFactory.setupProxy(l1CrossDomainMessenger, proxyAdmin, salt),
            l1StandardBridge: ResolvingProxyFactory.setupProxy(l1StandardBridge, proxyAdmin, salt),
            l1ERC721Bridge: ResolvingProxyFactory.setupProxy(l1ERC721Bridge, proxyAdmin, salt),
            optimismMintableERC20Factory: ResolvingProxyFactory.setupProxy(optimismMintableERC20Factory, proxyAdmin, salt)
        });
    }

    function calculateHashes(
        uint256 chainID,
        GenesisConfiguration memory genesisConfig,
        GasConfiguration memory gasConfig,
        address batcherAddress,
        DeployAddresses memory addresses
    ) internal view returns (Hashes memory) {
        bytes32 genesisL1Hash = blockhash(uint256(genesisConfig.l1Number));
        require(genesisL1Hash != bytes32(0), "DeployChain: genesis blockhash not available");

        bytes32 scalar =
            bytes32((uint256(0x01) << 248) | (uint256(gasConfig.blobbasefeeScalar) << 32) | gasConfig.basefeeScalar);

        bytes32 configHash = keccak256(
            abi.encodePacked(
                uint64(0), // version
                chainID,
                genesisL1Hash,
                genesisConfig.l2Hash,
                genesisConfig.l2Time,
                batcherAddress,
                scalar,
                gasConfig.gasLimit,
                addresses.optimismPortal,
                addresses.systemConfig
            )
        );

        bytes32 genesisOutputRoot = Hashing.hashOutputRootProof(
            Types.OutputRootProof({
                version: 0,
                stateRoot: genesisConfig.l2StateRoot,
                messagePasserStorageRoot: MESSAGE_PASSER_STORAGE_HASH,
                latestBlockhash: genesisConfig.l2Hash
            })
        );

        return Hashes({configHash: configHash, genesisOutputRoot: genesisOutputRoot});
    }

    function initializeProxies(
        GasConfiguration memory gasConfig,
        AddressConfiguration memory addressConfig,
        address batchInbox,
        Hashes memory hashes,
        DeployAddresses memory addresses
    ) internal {
        OutputOracle(addresses.l2OutputOracle).initialize(
            SystemConfigOwnable(addresses.systemConfig), hashes.configHash, hashes.genesisOutputRoot
        );

        Portal(payable(addresses.optimismPortal)).initialize(
            OutputOracle(addresses.l2OutputOracle),
            ISystemConfig(addresses.systemConfig),
            ISuperchainConfig(superchainConfig)
        );

        SystemConfig.Addresses memory systemAddresses = _createSystemAddresses(addresses, gasConfig.gasToken);

        SystemConfigOwnable(addresses.systemConfig).initialize({
            _basefeeScalar: gasConfig.basefeeScalar,
            _blobbasefeeScalar: gasConfig.blobbasefeeScalar,
            _batcherHash: bytes32(uint256(uint160(addressConfig.batcher))),
            _gasLimit: gasConfig.gasLimit,
            _unsafeBlockSigner: addressConfig.unsafeBlockSigner,
            _config: Constants.DEFAULT_RESOURCE_CONFIG(),
            _batchInbox: batchInbox,
            _proposer: addressConfig.proposer,
            _addresses: systemAddresses
        });

        L1CrossDomainMessenger(addresses.l1CrossDomainMessenger).initialize(
            ISuperchainConfig(superchainConfig),
            IOptimismPortal(payable(addresses.optimismPortal)),
            ISystemConfig(addresses.systemConfig)
        );

        L1StandardBridge(payable(addresses.l1StandardBridge)).initialize(
            ICrossDomainMessenger(addresses.l1CrossDomainMessenger),
            ISuperchainConfig(superchainConfig),
            ISystemConfig(addresses.systemConfig)
        );

        L1ERC721Bridge(addresses.l1ERC721Bridge).initialize(
            ICrossDomainMessenger(addresses.l1CrossDomainMessenger), ISuperchainConfig(superchainConfig)
        );

        OptimismMintableERC20Factory(addresses.optimismMintableERC20Factory).initialize(addresses.l1StandardBridge);
    }

    function _createSystemAddresses(DeployAddresses memory addresses, address gasToken)
        private
        pure
        returns (SystemConfig.Addresses memory)
    {
        return SystemConfig.Addresses({
            l1CrossDomainMessenger: addresses.l1CrossDomainMessenger,
            l1ERC721Bridge: addresses.l1ERC721Bridge,
            l1StandardBridge: addresses.l1StandardBridge,
            disputeGameFactory: address(0),
            optimismPortal: addresses.optimismPortal,
            optimismMintableERC20Factory: addresses.optimismMintableERC20Factory,
            gasPayingToken: gasToken
        });
    }
}
