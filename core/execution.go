// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	"math/big"

	"github.com/eth-classic/go-ethereum/common"
	"github.com/eth-classic/go-ethereum/core/vm"
	"github.com/eth-classic/go-ethereum/crypto"
)

var (
	callCreateDepthMax = 1024 // limit call/create stack
	errCallCreateDepth = fmt.Errorf("Max call depth exceeded (%d)", callCreateDepthMax)

	maxCodeSize            = 24576
	errMaxCodeSizeExceeded = fmt.Errorf("Max Code Size exceeded (%d)", maxCodeSize)
)

// Call executes within the given contract
func Call(env vm.Environment, caller vm.ContractRef, addr common.Address, input []byte, gas, gasPrice, value *big.Int) (ret []byte, err error) {

	// Depth check. Fail if we're trying to execute above the call depth limit
	if env.Depth() > callCreateDepthMax {
		caller.ReturnGas(gas, gasPrice)

		return nil, errCallCreateDepth
	}

	// Balance check. Fail if we're trying to transfer more than the available balance
	if !env.CanTransfer(caller.Address(), value) {
		caller.ReturnGas(gas, gasPrice)
		return nil,  ValueTransferErr("insufficient funds to transfer value. Req %v, has %v", value, env.Db().GetBalance(caller.Address()))
	}

	var (
		snapshotPreTransfer = env.SnapshotDatabase()
		to   vm.Account
		from = env.Db().GetAccount(caller.Address())
	)

	if !env.Db().Exist(addr) {
		//no account may change state from non-existent to existent-but-empty. Refund sender.
		if vm.PrecompiledAtlantis[(addr).Str()] == nil && env.RuleSet().IsAtlantis(env.BlockNumber()) && value.BitLen() == 0 {
			caller.ReturnGas(gas, gasPrice)
			return nil, nil
		}
		to = env.Db().CreateAccount(addr)
	} else {
		to = env.Db().GetAccount(addr)
	}

	env.Transfer(from, to, value)

	// initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := vm.NewContract(caller, to, value, gas, gasPrice)
	contract.SetCallCode(&addr, env.Db().GetCodeHash(addr), env.Db().GetCode(addr))
	defer contract.Finalise()

	ret, err = env.Vm().Run(contract, input, false)

	if err != nil {
		if err != vm.ErrRevert {
			contract.UseGas(contract.Gas)
		}
		env.RevertToSnapshot(snapshotPreTransfer)
	}
	return ret, err

}

// CallCode executes the given address' code as the given contract address
func CallCode(env vm.Environment, caller vm.ContractRef, addr common.Address, input []byte, gas, gasPrice, value *big.Int) (ret []byte, err error) {

	// Depth check. Fail if we're trying to execute above the call depth limit
	if env.Depth() > callCreateDepthMax {
		caller.ReturnGas(gas, gasPrice)

		return nil, errCallCreateDepth
	}

	// Balance check. Fail if we're trying to transfer more than the available balance
	if !env.CanTransfer(caller.Address(), value) {
		caller.ReturnGas(gas, gasPrice)
		return nil,  ValueTransferErr("insufficient funds to transfer value. Req %v, has %v", value, env.Db().GetBalance(caller.Address()))
	}

	var (
		snapshotPreTransfer = env.SnapshotDatabase()
		to   vm.Account
		from = env.Db().GetAccount(caller.Address())
	)

	if !env.Db().Exist(caller.Address()) {
		//no account may change state from non-existent to existent-but-empty. Refund sender.
		if vm.PrecompiledAtlantis[(caller.Address()).Str()] == nil && env.RuleSet().IsAtlantis(env.BlockNumber()) && value.BitLen() == 0 {
			caller.ReturnGas(gas, gasPrice)
			return nil, nil
		}
		to = env.Db().CreateAccount(caller.Address())
	} else {
		to = env.Db().GetAccount(caller.Address())
	}

	env.Transfer(from, to, value)

	// initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := vm.NewContract(caller, to, value, gas, gasPrice)
	contract.SetCallCode(&addr, env.Db().GetCodeHash(addr), env.Db().GetCode(addr))
	defer contract.Finalise()

	ret, err = env.Vm().Run(contract, input, false)

	if err != nil {
		if err != vm.ErrRevert {
			contract.UseGas(contract.Gas)
		}
		env.RevertToSnapshot(snapshotPreTransfer)
	}
	return ret, err
}

// DelegateCall is equivalent to CallCode except that sender and value propagates from parent scope to child scope
func DelegateCall(env vm.Environment, caller vm.ContractRef, addr common.Address, input []byte, gas, gasPrice *big.Int) (ret []byte, err error) {

	// Depth check execution. Fail if we're trying to execute above the limit.
	if env.Depth() > callCreateDepthMax {
		caller.ReturnGas(gas, gasPrice)

		return nil, errCallCreateDepth
	}

	var (
		to       vm.Account
		snapshot = env.SnapshotDatabase()
	)
	if !env.Db().Exist(caller.Address()) {
		to = env.Db().CreateAccount(caller.Address())
	} else {
		to = env.Db().GetAccount(caller.Address())
	}

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := vm.NewContract(caller, to, caller.Value(), gas, gasPrice).AsDelegate()
	contract.SetCallCode(&addr, env.Db().GetCodeHash(addr), env.Db().GetCode(addr))
	defer contract.Finalise()

	ret, err = env.Vm().Run(contract, input, false)

	if err != nil {
		env.RevertToSnapshot(snapshot)
		if err != vm.ErrRevert {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, err
}

// StaticCall executes within the given contract and throws exception if state is attempted to be changed
func StaticCall(env vm.Environment, caller vm.ContractRef, addr common.Address, input []byte, gas, gasPrice *big.Int) (ret []byte, err error) {
	
	// Depth check execution. Fail if we're trying to execute above the limit.
	if env.Depth() > callCreateDepthMax {
		caller.ReturnGas(gas, gasPrice)

		return nil, errCallCreateDepth
	}

	var (
		to       vm.Account
		snapshot = env.SnapshotDatabase()
	)
	if !env.Db().Exist(addr) {
		to = env.Db().CreateAccount(addr)
	} else {
		to = env.Db().GetAccount(addr)
	}

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := vm.NewContract(caller, to, new(big.Int), gas, gasPrice)
	contract.SetCallCode(&addr, env.Db().GetCodeHash(addr), env.Db().GetCode(addr))
	defer contract.Finalise()

	//STATICCALL touch delete
	//see https://github.com/ethereum/go-ethereum/pull/18187
	env.Db().AddBalance(addr, big.NewInt(0))

	ret, err = env.Vm().Run(contract, input, true)

	if err != nil {
		env.RevertToSnapshot(snapshot)
		if err != vm.ErrRevert {
			contract.UseGas(contract.Gas)
		}
	}
	return ret, err
}

// Create creates a new contract with the given code
func Create(env vm.Environment, caller vm.ContractRef, code []byte, gas, gasPrice, value *big.Int) (ret []byte, address common.Address, err error) {

	// Depth check execution. Fail if we're trying to execute above the limit.
	if env.Depth() > callCreateDepthMax {
		caller.ReturnGas(gas, gasPrice)

		return nil, common.Address{}, errCallCreateDepth
	}
	if !env.CanTransfer(caller.Address(), value) {
		caller.ReturnGas(gas, gasPrice)

		return nil, common.Address{}, ErrInsufficientFunds
	}
	nonce := env.Db().GetNonce(caller.Address())
	env.Db().SetNonce(caller.Address(), nonce+1)

	// Create a new account on the state
	snapshot := env.SnapshotDatabase()

	addr := crypto.CreateAddress(caller.Address(), nonce)

	//Create account with address
	to := env.Db().CreateAccount(addr)

	if env.RuleSet().IsAtlantis(env.BlockNumber()) {
		env.Db().SetNonce(addr, 1)
	}

	env.Transfer(env.Db().GetAccount(caller.Address()), to, value)

	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := vm.NewContract(caller, to, value, gas, gasPrice)
	contract.SetCallCode(nil, crypto.Keccak256Hash(code), code)
	defer contract.Finalise()

	ret, err = env.Vm().Run(contract, nil, false)

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := len(ret) > maxCodeSize && env.RuleSet().IsAtlantis(env.BlockNumber())
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := big.NewInt(int64(len(ret)))
		createDataGas.Mul(createDataGas, big.NewInt(200))
		if contract.UseGas(createDataGas) {
			env.Db().SetCode(addr, ret)
		} else {
			err = vm.CodeStoreOutOfGasError
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && (env.RuleSet().IsHomestead(env.BlockNumber()) || err != vm.CodeStoreOutOfGasError)) {
		env.RevertToSnapshot(snapshot)
		if err != vm.ErrRevert {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}

	//if there's an error we return nothing
	if err != nil && err != vm.ErrRevert {
		return nil, addr, err
	}
	return ret, addr, err
}

// generic transfer method
func Transfer(from, to vm.Account, amount *big.Int) {
	from.SubBalance(amount)
	to.AddBalance(amount)
}
