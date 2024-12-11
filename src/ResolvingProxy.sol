// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

// The ProxyAdmin contract conforms to this interface.
interface IResolver {
    function getProxyImplementation(address _proxy) external view returns (address);
}

/// @notice ResolvingProxy is a modification of the op-stack Proxy contract that allows for
///         proxying a proxy. This is useful to have a central upgradable proxy that this
///         contract can point to, but also support detaching this proxy so it can have its
///         own implementation.
/// @dev Only proxies that are owned by a ProxyAdmin can be proxied by this contract,
///      because it calls getProxyImplementation() on the ProxyAdmin to retrieve the
///      implementation address.
/// @dev This contract is based on the EIP-1967 transparent proxy standard. It is slightly
///      simplified in that it doesn't emit logs for implementation and admin changes. This
///      is to simplify the assembly implementation provided in ResolvingProxyFactory.
contract ResolvingProxy {
    /// @notice The storage slot that holds the address of a proxy implementation.
    /// @dev `bytes32(uint256(keccak256('eip1967.proxy.implementation')) - 1)`
    bytes32 internal constant IMPLEMENTATION_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;

    /// @notice The storage slot that holds the address of the owner.
    /// @dev `bytes32(uint256(keccak256('eip1967.proxy.admin')) - 1)`
    bytes32 internal constant ADMIN_SLOT = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    /// @notice A modifier that reverts if not called by the owner or by address(0) to allow
    ///         eth_call to interact with this proxy without needing to use low-level storage
    ///         inspection. We assume that nobody is able to trigger calls from address(0) during
    ///         normal EVM execution.
    modifier proxyCallIfNotAdmin() {
        if (msg.sender == _getAdmin() || msg.sender == address(0)) {
            _;
        } else {
            // This WILL halt the call frame on completion.
            _doProxyCall();
        }
    }

    /// @notice Sets the initial admin during contract deployment. Admin address is stored at the
    ///         EIP-1967 admin storage slot so that accidental storage collision with the
    ///         implementation is not possible.
    /// @param _admin Address of the initial contract admin. Admin has the ability to access the
    ///               transparent proxy interface.
    constructor(address _admin) {
        _setAdmin(_admin);
    }

    receive() external payable {
        // Proxy call by default.
        _doProxyCall();
    }

    fallback() external payable {
        // Proxy call by default.
        _doProxyCall();
    }

    /// @notice Gets the owner of the proxy contract.
    /// @return Owner address.
    function admin() external virtual proxyCallIfNotAdmin returns (address) {
        return _getAdmin();
    }

    /// @notice Changes the owner of the proxy contract. Only callable by the owner.
    /// @param _admin New owner of the proxy contract.
    function changeAdmin(address _admin) external virtual proxyCallIfNotAdmin {
        _setAdmin(_admin);
    }

    //// @notice Queries the implementation address.
    /// @return Implementation address.
    function implementation() external virtual proxyCallIfNotAdmin returns (address) {
        return _getImplementation();
    }

    /// @notice Set the implementation contract address. The code at the given address will execute
    ///         when this contract is called.
    /// @param _implementation Address of the implementation contract.
    function upgradeTo(address _implementation) external virtual proxyCallIfNotAdmin {
        _setImplementation(_implementation);
    }

    /// @notice Set the implementation and call a function in a single transaction. Useful to ensure
    ///         atomic execution of initialization-based upgrades.
    /// @param _implementation Address of the implementation contract.
    /// @param _data           Calldata to delegatecall the new implementation with.
    function upgradeToAndCall(address _implementation, bytes calldata _data)
        external
        payable
        proxyCallIfNotAdmin
        returns (bytes memory)
    {
        _setImplementation(_implementation);
        _implementation = _resolveImplementation();
        (bool success, bytes memory returndata) = _implementation.delegatecall(_data);
        require(success, "Proxy: delegatecall to new implementation contract failed");
        return returndata;
    }

    function _getImplementation() internal view returns (address) {
        address impl;
        bytes32 proxyImplementation = IMPLEMENTATION_SLOT;
        assembly {
            impl := sload(proxyImplementation)
        }
        return impl;
    }

    function _setImplementation(address _implementation) internal {
        bytes32 proxyImplementation = IMPLEMENTATION_SLOT;
        assembly {
            sstore(proxyImplementation, _implementation)
        }
    }

    function _getAdmin() internal view returns (address) {
        address owner;
        bytes32 proxyOwner = ADMIN_SLOT;
        assembly {
            owner := sload(proxyOwner)
        }
        return owner;
    }

    function _setAdmin(address _admin) internal {
        bytes32 proxyOwner = ADMIN_SLOT;
        assembly {
            sstore(proxyOwner, _admin)
        }
    }

    function _doProxyCall() internal {
        address impl = _resolveImplementation();
        assembly {
            calldatacopy(0x0, 0x0, calldatasize())
            let success := delegatecall(gas(), impl, 0x0, calldatasize(), 0x0, 0x0)
            returndatacopy(0x0, 0x0, returndatasize())
            if iszero(success) { revert(0x0, returndatasize()) }
            return(0x0, returndatasize())
        }
    }

    function _resolveImplementation() internal view returns (address) {
        address proxy = _getImplementation();
        require(proxy != address(0));
        bytes memory data = abi.encodeCall(IResolver.getProxyImplementation, (proxy));
        (bool success, bytes memory returndata) = _getAdmin().staticcall(data);
        if (success && returndata.length == 0x20) {
            return abi.decode(returndata, (address));
        }
        return proxy;
    }
}
