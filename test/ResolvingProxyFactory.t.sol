// SPDX-License-Identifier: MIT
pragma solidity ^0.8.15;

import {Vm} from "forge-std/Vm.sol";
import {Test, console} from "forge-std/Test.sol";
import {ResolvingProxyFactory} from "../src/ResolvingProxyFactory.sol";
import {ProxyAdmin} from "@eth-optimism-bedrock/src/universal/ProxyAdmin.sol";
import {Proxy} from "@eth-optimism-bedrock/src/universal/Proxy.sol";

abstract contract Implementation {
    string internal value;

    function set(string memory _value) external {
        value = _value;
    }
}

contract Implementation1 is Implementation {
    function get() external view returns (string memory) {
        return string.concat("Hello ", value);
    }
}

contract Implementation2 is Implementation {
    function get() external view returns (string memory) {
        return string.concat("Hi ", value);
    }
}

contract ResolvingProxyFactoryTest is Test {
    Implementation1 public implementation1;
    Implementation2 public implementation2;
    ProxyAdmin public admin;
    Proxy public proxy;
    address public resolvingProxy;

    function setUp() public {
        implementation1 = new Implementation1();
        implementation2 = new Implementation2();
        admin = new ProxyAdmin(address(this));
        proxy = new Proxy(address(admin));
        admin.upgrade(payable(address(proxy)), address(implementation1));
        resolvingProxy = ResolvingProxyFactory.setupProxy(address(admin), 0x00);
        admin.upgrade(payable(address(resolvingProxy)), address(proxy));
        Implementation1(resolvingProxy).set("world");
    }

    function test_setupProxy() public view {
        string memory value = Implementation1(resolvingProxy).get();
        assertEq(value, "Hello world");
    }

    function test_implementation() public view {
        address implementation = admin.getProxyImplementation(payable(address(resolvingProxy)));
        assertEq(implementation, address(proxy));
    }

    function test_admin() public view {
        address _admin = admin.getProxyAdmin(payable(address(resolvingProxy)));
        assertEq(_admin, address(admin));
    }

    function test_upgradeTo() public {
        admin.upgrade(payable(address(resolvingProxy)), address(implementation2));
        Implementation1(resolvingProxy).set("alice");
        string memory value = Implementation1(resolvingProxy).get();
        assertEq(value, "Hi alice");

        admin.upgrade(payable(address(resolvingProxy)), address(proxy));
        value = Implementation1(resolvingProxy).get();
        assertEq(value, "Hello alice");
    }

    function test_upgradeToAndCall() public {
        bytes memory data = abi.encodeCall(Implementation.set, ("alice"));
        admin.upgradeAndCall(payable(address(resolvingProxy)), address(implementation2), data);
        string memory value = Implementation1(resolvingProxy).get();
        assertEq(value, "Hi alice");
    }

    function test_changeAdmin() public {
        ProxyAdmin newAdmin = new ProxyAdmin(address(this));
        admin.changeProxyAdmin(payable(address(resolvingProxy)), address(newAdmin));
        newAdmin.upgrade(payable(address(resolvingProxy)), address(implementation2));
        string memory value = Implementation1(resolvingProxy).get();
        assertEq(value, "Hi world");
    }

    function test_proxyAddress() public view {
        address predicted = ResolvingProxyFactory.proxyAddress(address(admin), 0x00);
        assertEq(predicted, resolvingProxy);
    }
}
