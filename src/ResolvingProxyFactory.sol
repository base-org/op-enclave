// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import {ResolvingProxy} from "./ResolvingProxy.sol";

/// @title ResolvingProxyFactory
/// @notice ResolvingProxyFactory is a factory contract that creates ResolvingProxy instances.
/// @dev The setupProxy / proxyAddress functions provide a smaller assembly-based ResolvingProxy
///      implementation that is more gas efficient to deploy and operate than the solidity
///      ResolvingProxy implementation.
library ResolvingProxyFactory {
    function setupProxy(address proxy, address admin, bytes32 salt) internal returns (address instance) {
        /// @solidity memory-safe-assembly
        assembly {
            let ptr := mload(0x40)
            mstore(ptr, shl(0xC0, 0x600661011c565b73))
            mstore(add(ptr, 0x8), shl(0x60, proxy))
            mstore(add(ptr, 0x1c), shl(0xE8, 0x905573))
            mstore(add(ptr, 0x1f), shl(0x60, admin))
            mstore(add(ptr, 0x33), 0x905561012580603f5f395ff35f365f600860dd565b8054909180548033143315)
            mstore(add(ptr, 0x53), 0x171560545760045f5f375f5160e01c8063f851a4401460a25780635c60da1b14)
            mstore(add(ptr, 0x73), 0x609f5780638f2839701460af5780633659cfe61460ac57634f1ef2861460aa57)
            mstore(add(ptr, 0x93), 0x5b63204e1c7a60e01b5f52826004525f5f60245f845afa3d5f5f3e3d60201416)
            mstore(add(ptr, 0xb3), 0x805f510290158402015f875f89895f375f935af43d5f893d60205260205f523e)
            mstore(add(ptr, 0xd3), 0x5f3d890191609d57fd5bf35b50505b505f5260205ff35b5f5b93915b50506020)
            mstore(add(ptr, 0xf3), 0x60045f375f518091559160d957903333602060445f375f519560649550506040)
            mstore(add(ptr, 0x113), 0x96506054565b5f5ff35b7f360894a13ba1a3210667c828492db98dca3e2076cc)
            mstore(add(ptr, 0x133), 0x3735a920a3ca505d382bbc7fb53127684a568b3173ae13b9f8a6016e243e63b6)
            mstore(add(ptr, 0x153), shl(0x90, 0xe8ee1178d6a717850b5d61039156))
            instance := create2(0, ptr, 0x161, salt)
        }
        require(instance != address(0), "Proxy: create2 failed");
    }

    function proxyAddress(address proxy, address admin, bytes32 salt) internal view returns (address predicted) {
        address deployer = address(this);
        /// @solidity memory-safe-assembly
        assembly {
            let ptr := mload(0x40)
            mstore(ptr, shl(0xC0, 0x600661011c565b73))
            mstore(add(ptr, 0x8), shl(0x60, proxy))
            mstore(add(ptr, 0x1c), shl(0xE8, 0x905573))
            mstore(add(ptr, 0x1f), shl(0x60, admin))
            mstore(add(ptr, 0x33), 0x905561012580603f5f395ff35f365f600860dd565b8054909180548033143315)
            mstore(add(ptr, 0x53), 0x171560545760045f5f375f5160e01c8063f851a4401460a25780635c60da1b14)
            mstore(add(ptr, 0x73), 0x609f5780638f2839701460af5780633659cfe61460ac57634f1ef2861460aa57)
            mstore(add(ptr, 0x93), 0x5b63204e1c7a60e01b5f52826004525f5f60245f845afa3d5f5f3e3d60201416)
            mstore(add(ptr, 0xb3), 0x805f510290158402015f875f89895f375f935af43d5f893d60205260205f523e)
            mstore(add(ptr, 0xd3), 0x5f3d890191609d57fd5bf35b50505b505f5260205ff35b5f5b93915b50506020)
            mstore(add(ptr, 0xf3), 0x60045f375f518091559160d957903333602060445f375f519560649550506040)
            mstore(add(ptr, 0x113), 0x96506054565b5f5ff35b7f360894a13ba1a3210667c828492db98dca3e2076cc)
            mstore(add(ptr, 0x133), 0x3735a920a3ca505d382bbc7fb53127684a568b3173ae13b9f8a6016e243e63b6)
            mstore(add(ptr, 0x153), shl(0x88, 0xe8ee1178d6a717850b5d61039156ff))
            mstore(add(ptr, 0x162), shl(0x60, deployer))
            mstore(add(ptr, 0x176), salt)
            mstore(add(ptr, 0x196), keccak256(ptr, 0x161))
            predicted := keccak256(add(ptr, 0x161), 0x55)
        }
    }

    function setupExpensiveProxy(address proxy, address admin, bytes32 salt) internal returns (address instance) {
        return address(new ResolvingProxy{salt: salt}(proxy, admin));
    }

    function expensiveProxyAddress(address proxy, address admin, bytes32 salt) internal view returns (address predicted) {
        bytes memory bytecode = abi.encodePacked(type(ResolvingProxy).creationCode, abi.encode(proxy, admin));
        bytes32 hash = keccak256(abi.encodePacked(bytes1(0xff), address(this), salt, keccak256(bytecode)));
        return address(uint160(uint256(hash)));
    }
}
