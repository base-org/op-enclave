// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

// Contracts
import {Initializable} from "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import {SystemConfigGlobal} from "./SystemConfigGlobal.sol";
import {SystemConfigOwnable} from "./SystemConfigOwnable.sol";

// Libraries
import {Types} from "@eth-optimism-bedrock/src/libraries/Types.sol";
import {Constants} from "@eth-optimism-bedrock/src/libraries/Constants.sol";
import {ECDSA} from "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

// Interfaces
import {ISemver} from "@eth-optimism-bedrock/src/universal/interfaces/ISemver.sol";

/// @custom:proxied true
/// @title OutputOracle
/// @notice The OutputOracle contains an array of L2 state outputs, where each output is a
///         commitment to the state of the L2 chain. Other contracts like the Portal use
///         these outputs to verify information about the state of L2.
///
///         Copy of https://github.com/ethereum-optimism/optimism/blob/0cf2a11611d329b9c6d53fd6473dbbc44f68e94c/packages/contracts-bedrock/src/L1/L2OutputOracle.sol
///         with the following modifications
///          - Limited number of outputs to store in l2Outputs
///          - Removed finalizationPeriod
///          - Additional signature validation for proposeL2Output
///          - Allow outputs to be proposed at any time
contract OutputOracle is Initializable, ISemver {
    /// @notice Global system config.
    SystemConfigGlobal public immutable systemConfigGlobal;

    /// @notice Maximum number of outputs to store in l2Outputs.
    uint256 public immutable maxOutputCount;

    /// @notice Per chain system config contract.
    SystemConfigOwnable public systemConfig;

    /// @notice Hash of the serialized chain configuration.
    bytes32 public configHash;

    /// @notice An array of L2 output proposals.
    Types.OutputProposal[] internal l2Outputs;

    /// @notice Whether or not TEE proofs are enabled.
    bool public proofsEnabled;

    /// @notice Pointer inside l2Outputs to the latest submitted output.
    uint256 public latestOutputIndex;

    /// @notice Emitted when an output is proposed.
    /// @param outputRoot    The output root.
    /// @param l2OutputIndex The index of the output in the l2Outputs array.
    /// @param l2BlockNumber The L2 block number of the output root.
    /// @param l1Timestamp   The L1 timestamp when proposed.
    event OutputProposed(
        bytes32 indexed outputRoot, uint256 indexed l2OutputIndex, uint256 indexed l2BlockNumber, uint256 l1Timestamp
    );

    /// @notice Semantic version.
    /// @custom:semver 1.8.1-beta.1
    string public constant version = "1.8.1-beta.1";

    /// @notice Constructs the OutputOracle contract. Initializes variables to the same values as
    ///         in the getting-started config.
    /// @param _systemConfigGlobal  The global system config contract.
    /// @param _maxOutputCount      The maximum number of outputs stored by this contract.
    constructor(SystemConfigGlobal _systemConfigGlobal, uint256 _maxOutputCount) {
        systemConfigGlobal = _systemConfigGlobal;
        maxOutputCount = _maxOutputCount;
        initialize(SystemConfigOwnable(address(0)), 0, 0, false);
    }

    /// @notice Initializer.
    function initialize(
        SystemConfigOwnable _systemConfig,
        bytes32 _configHash,
        bytes32 _genesisOutputRoot,
        bool _proofsEnabled
    ) public initializer {
        systemConfig = _systemConfig;
        configHash = _configHash;
        l2Outputs.push(
            Types.OutputProposal({
                outputRoot: _genesisOutputRoot,
                timestamp: uint128(block.timestamp),
                l2BlockNumber: uint128(0)
            })
        );
        proofsEnabled = _proofsEnabled;
    }

    /// @notice Returns the proposer address.
    function proposer() public view returns (address) {
        address _proposer = address(systemConfig) != address(0) ? systemConfig.proposer() : address(0);
        return _proposer != address(0) ? _proposer : systemConfigGlobal.proposer();
    }

    /// @notice Enables TEE proofs for output proposals. This is a one-way function.
    function enableProofs() external {
        require(msg.sender == proposer(), "OutputOracle: only the proposer address can enable proofs");
        proofsEnabled = true;
    }

    /// @notice Accepts an outputRoot of the corresponding L2 block.
    /// @param _outputRoot    The L2 output of the checkpoint block.
    /// @param _l2BlockNumber The L2 block number that resulted in _outputRoot.
    function proposeL2Output(
        bytes32 _outputRoot,
        uint256 _l2BlockNumber,
        uint256 _l1BlockNumber,
        bytes calldata _signature
    ) external {
        require(msg.sender == proposer(), "OutputOracle: only the proposer address can propose new outputs");

        require(
            _l2BlockNumber > latestBlockNumber(),
            "OutputOracle: block number must be greater than previously proposed block number"
        );

        require(_outputRoot != bytes32(0), "OutputOracle: L2 output proposal cannot be the zero hash");

        if (proofsEnabled) {
            bytes32 _blockHash = blockhash(_l1BlockNumber);
            require(_blockHash != bytes32(0), "OutputOracle: blockhash not available");

            bytes32 previousOutputRoot = l2Outputs[latestOutputIndex].outputRoot;
            address signer = ECDSA.recover(
                keccak256(abi.encodePacked(configHash, _blockHash, previousOutputRoot, _outputRoot)), _signature
            );
            require(systemConfigGlobal.validSigners(signer), "OutputOracle: invalid signature");
        }

        latestOutputIndex = nextOutputIndex();
        emit OutputProposed(_outputRoot, latestOutputIndex, _l2BlockNumber, block.timestamp);

        Types.OutputProposal memory op = Types.OutputProposal({
            outputRoot: _outputRoot,
            timestamp: uint128(block.timestamp),
            l2BlockNumber: uint128(_l2BlockNumber)
        });
        if (l2Outputs.length < maxOutputCount) {
            l2Outputs.push(op);
        } else {
            l2Outputs[latestOutputIndex] = op;
        }
    }

    /// @notice Returns an output by index. Needed to return a struct instead of a tuple.
    /// @param _l2OutputIndex Index of the output to return.
    /// @return The output at the given index.
    function getL2Output(uint256 _l2OutputIndex) external view returns (Types.OutputProposal memory) {
        return l2Outputs[_l2OutputIndex];
    }

    /// @notice Returns the index of the L2 output that checkpoints a given L2 block number.
    ///         Uses a binary search to find the first output greater than or equal to the given
    ///         block.
    ///         This is provided for backwards compatibility with the op-stack. However, because
    ///         this contract is modified to only store a fixed number of outputs, it is recommended
    ///         to use the `latestOutputIndex` rather than the index returned from this function,
    ///         as it will stay valid for longer.
    ///         If your withdrawal is old, the output at the index selected by this function may be
    ///         overwritten by the next proposal.
    /// @param _l2BlockNumber L2 block number to find a checkpoint for.
    /// @return Index of the first checkpoint that commits to the given L2 block number.
    function getL2OutputIndexAfter(uint256 _l2BlockNumber) public view returns (uint256) {
        // Make sure an output for this block number has actually been proposed.
        require(
            _l2BlockNumber <= latestBlockNumber(),
            "L2OutputOracle: cannot get output for a block that has not been proposed"
        );

        // Find the output via binary search, guaranteed to exist.
        uint256 offset = nextOutputIndex();
        uint256 lo = offset;
        uint256 hi = offset + l2Outputs.length;
        while (lo < hi) {
            uint256 mid = (lo + hi) / 2;
            if (l2Outputs[mid % l2Outputs.length].l2BlockNumber < _l2BlockNumber) {
                lo = mid + 1;
            } else {
                hi = mid;
            }
        }

        return lo % l2Outputs.length;
    }

    /// @notice Returns the L2 output proposal that checkpoints a given L2 block number.
    ///         Uses a binary search to find the first output greater than or equal to the given
    ///         block.
    /// @param _l2BlockNumber L2 block number to find a checkpoint for.
    /// @return First checkpoint that commits to the given L2 block number.
    function getL2OutputAfter(uint256 _l2BlockNumber) external view returns (Types.OutputProposal memory) {
        return l2Outputs[getL2OutputIndexAfter(_l2BlockNumber)];
    }

    /// @notice Returns the index of the next output to be proposed.
    /// @return The index of the next output to be proposed.
    function nextOutputIndex() public view returns (uint256) {
        return (latestOutputIndex + 1) % maxOutputCount;
    }

    /// @notice Returns the block number of the latest submitted L2 output proposal.
    ///         If no proposals been submitted yet then this function will return 0.
    /// @return Latest submitted L2 block number.
    function latestBlockNumber() public view returns (uint256) {
        return l2Outputs.length == 0 ? 0 : l2Outputs[latestOutputIndex].l2BlockNumber;
    }

    /// @notice Returns the latest submitted L2 output proposal.
    /// @return Latest submitted L2 output proposal.
    function latestL2Output() public view returns (Types.OutputProposal memory) {
        return l2Outputs[latestOutputIndex];
    }
}
