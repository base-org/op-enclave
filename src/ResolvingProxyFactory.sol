// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import {ResolvingProxy} from "./ResolvingProxy.sol";

/// @title ResolvingProxyFactory
/// @notice ResolvingProxyFactory is a factory contract that creates ResolvingProxy instances.
/// @dev The setupProxy / proxyAddress functions provide a smaller assembly-based ResolvingProxy
///      implementation that is more gas efficient to deploy and operate than the solidity
///      ResolvingProxy implementation.
library ResolvingProxyFactory {
    function setupProxy(address admin, bytes32 salt) internal returns (address instance) {
        /// @solidity memory-safe-assembly
        assembly {
            let ptr := mload(0x40)
            mstore(ptr, shl(0xC0, 0x600661010e565b73))
            mstore(add(ptr, 0x8), shl(0x60, admin))
            mstore(add(ptr, 0x1c), 0x90915561012a8060295f395ff35f365f600860e5565b80549091805480331433)
            mstore(add(ptr, 0x3c), 0x15171560545760045f5f375f5160e01c8063f851a4401460a75780635c60da1b)
            mstore(add(ptr, 0x5c), 0x1460a45780638f2839701460b45780633659cfe61460b157634f1ef2861460af)
            mstore(add(ptr, 0x7c), 0x575b63204e1c7a60e01b5f5282801560e2576004525f5f60245f845afa3d5f5f)
            mstore(add(ptr, 0x9c), 0x3e3d60201416805f510290158402015f875f89895f375f935af43d5f893d6020)
            mstore(add(ptr, 0xbc), 0x5260205f523e5f3d89019160a257fd5bf35b50505b505f5260205ff35b5f5b93)
            mstore(add(ptr, 0xdc), 0x915b5050602060045f375f518091559160de57903333602060445f375f519560)
            mstore(add(ptr, 0xfc), 0x64955050604096506054565b5f5ff35b5ffd5b7f360894a13ba1a3210667c828)
            mstore(add(ptr, 0x11c), 0x492db98dca3e2076cc3735a920a3ca505d382bbc7fb53127684a568b3173ae13)
            mstore(add(ptr, 0x13c), shl(0x48, 0xb9f8a6016e243e63b6e8ee1178d6a717850b5d61039156))
            instance := create2(0, ptr, 0x153, salt)
        }
        require(instance != address(0), "Proxy: create2 failed");
    }

    function proxyAddress(address admin, bytes32 salt) internal view returns (address predicted) {
        address deployer = address(this);
        /// @solidity memory-safe-assembly
        assembly {
            let ptr := mload(0x40)
            mstore(ptr, shl(0xC0, 0x600661010e565b73))
            mstore(add(ptr, 0x8), shl(0x60, admin))
            mstore(add(ptr, 0x1c), 0x90915561012a8060295f395ff35f365f600860e5565b80549091805480331433)
            mstore(add(ptr, 0x3c), 0x15171560545760045f5f375f5160e01c8063f851a4401460a75780635c60da1b)
            mstore(add(ptr, 0x5c), 0x1460a45780638f2839701460b45780633659cfe61460b157634f1ef2861460af)
            mstore(add(ptr, 0x7c), 0x575b63204e1c7a60e01b5f5282801560e2576004525f5f60245f845afa3d5f5f)
            mstore(add(ptr, 0x9c), 0x3e3d60201416805f510290158402015f875f89895f375f935af43d5f893d6020)
            mstore(add(ptr, 0xbc), 0x5260205f523e5f3d89019160a257fd5bf35b50505b505f5260205ff35b5f5b93)
            mstore(add(ptr, 0xdc), 0x915b5050602060045f375f518091559160de57903333602060445f375f519560)
            mstore(add(ptr, 0xfc), 0x64955050604096506054565b5f5ff35b5ffd5b7f360894a13ba1a3210667c828)
            mstore(add(ptr, 0x11c), 0x492db98dca3e2076cc3735a920a3ca505d382bbc7fb53127684a568b3173ae13)
            mstore(add(ptr, 0x13c), shl(0x40, 0xb9f8a6016e243e63b6e8ee1178d6a717850b5d61039156ff))
            mstore(add(ptr, 0x154), shl(0x60, deployer))
            mstore(add(ptr, 0x168), salt)
            mstore(add(ptr, 0x188), keccak256(ptr, 0x153))
            predicted := keccak256(add(ptr, 0x153), 0x55)
        }
    }

    function setupExpensiveProxy(address admin, bytes32 salt) internal returns (address instance) {
        return address(new ResolvingProxy{salt: salt}(admin));
    }

    function expensiveProxyAddress(address admin, bytes32 salt)
        internal
        view
        returns (address predicted)
    {
        bytes memory bytecode = abi.encodePacked(type(ResolvingProxy).creationCode, abi.encode(admin));
        bytes32 hash = keccak256(abi.encodePacked(bytes1(0xff), address(this), salt, keccak256(bytecode)));
        return address(uint160(uint256(hash)));
    }
}
