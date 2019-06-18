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
	"errors"
	"fmt"
	"math/big"

	"github.com/eth-classic/go-ethereum/common"
	"github.com/eth-classic/go-ethereum/common/hexutil"
	"github.com/eth-classic/go-ethereum/core/vm"
	"github.com/eth-classic/go-ethereum/crypto"
)

var (
	callCreateDepthMax = 1024 // limit call/create stack
	errCallCreateDepth = fmt.Errorf("Max call depth exceeded (%d)", callCreateDepthMax)

	maxCodeSize            = 24576
	errMaxCodeSizeExceeded = fmt.Errorf("Max Code Size exceeded (%d)", maxCodeSize)

	ErrCodeStoreOutOfGas        = errors.New("contract creation code storage out of gas")
<<<<<<< HEAD
=======
	errExecutionReverted        = errors.New("evm: execution reverted")
>>>>>>> Implemented create function for Create2 contracts
	ErrContractAddressCollision = errors.New("contract address collision")
)

// Call executes within the given contract
func Call(env vm.Environment, caller vm.ContractRef, addr common.Address, input []byte, gas, gasPrice, value *big.Int) (ret []byte, err error) {
	ret, _, err = exec(env, caller, &addr, &addr, env.Db().GetCodeHash(addr), input, env.Db().GetCode(addr), gas, gasPrice, value, false)
	return ret, err
}

// CallCode executes the given address' code as the given contract address
func CallCode(env vm.Environment, caller vm.ContractRef, addr common.Address, input []byte, gas, gasPrice, value *big.Int) (ret []byte, err error) {
	callerAddr := caller.Address()
	ret, _, err = exec(env, caller, &callerAddr, &addr, env.Db().GetCodeHash(addr), input, env.Db().GetCode(addr), gas, gasPrice, value, false)
	return ret, err
}

// DelegateCall is equivalent to CallCode except that sender and value propagates from parent scope to child scope
func DelegateCall(env vm.Environment, caller vm.ContractRef, addr common.Address, input []byte, gas, gasPrice *big.Int) (ret []byte, err error) {
	callerAddr := caller.Address()
	originAddr := env.Origin()
	callerValue := caller.Value()
	ret, _, err = execDelegateCall(env, caller, &originAddr, &callerAddr, &addr, env.Db().GetCodeHash(addr), input, env.Db().GetCode(addr), gas, gasPrice, callerValue)
	return ret, err
}

// StaticCall executes within the given contract and throws exception if state is attempted to be changed
func StaticCall(env vm.Environment, caller vm.ContractRef, addr common.Address, input []byte, gas, gasPrice *big.Int) (ret []byte, err error) {
	ret, _, err = exec(env, caller, &addr, &addr, env.Db().GetCodeHash(addr), input, env.Db().GetCode(addr), gas, gasPrice, new(big.Int), true)
	return ret, err
}

// Create creates a new contract with the given code
func Create(env vm.Environment, caller vm.ContractRef, code []byte, gas, gasPrice, value *big.Int) (ret []byte, address common.Address, err error) {
	ret, address, err = exec(env, caller, nil, nil, crypto.Keccak256Hash(code), nil, code, gas, gasPrice, value, false)
	// Here we get an error if we run into maximum stack depth,
	// See: https://github.com/ethereum/yellowpaper/pull/131
	// and YP definitions for CREATE

	//if there's an error we return nothing
	if err != nil && err != vm.ErrRevert {
		return nil, address, err
	}
	return ret, address, err
}

//They're Not using Sender nonce hash?

// Create2 creates a new contract with the given code
func Create2(env vm.Environment, caller vm.ContractRef, code []byte, gas, gasPrice, salt, value *big.Int) (ret []byte, address common.Address, err error) {
	addr := crypto.CreateAddress2(caller.Address(), common.BigToHash(salt).Bytes(), crypto.Keccak256(code))
<<<<<<< HEAD
	ret, address, err = create(env, caller, &addr, nil, crypto.Keccak256Hash(code), nil, code, gas, gasPrice, value, false)
=======
	fmt.Println("BLAH: " + hexutil.Encode(addr.Bytes()))
	// addr := crypto.CreateAddress2(caller.Address(), common.BigToHash(salt).Bytes(), crypto.Keccak256(code)[31:])
	ret, address, err = create(env, caller, &addr, nil, crypto.Keccak256Hash(code), nil, code, gas, gasPrice, value, false)
	fmt.Println("VALUE: ", value)
	// ret, address, err = exec(env, caller, &addr, nil, crypto.Keccak256Hash(code), nil, code, gas, gasPrice, value, false)
	fmt.Println("RESULTADDR: ", hexutil.Encode(addr.Bytes()))
	fmt.Println("ERROR: ", err)
	fmt.Println("CODE: ", ret)

	/*
		fmt.Println("oldAddress: ", hexutil.Encode(caller.Address().Bytes()))
		fmt.Println("SALT: " + salt.String())
		fmt.Println("Keccak256(init_code): ", hexutil.Encode(crypto.Keccak256Hash(code).Bytes()))
		// common.BigToHash()
		// hash := common.BytesToHash(crypto.aaaKeccak256([]byte{0xff}, caller.Address().Bytes(), common.BigToHash(salt).Bytes(), crypto.Keccak256Hash(code).Bytes())[12:])
		hash := common.BytesToHash(Keccak256([]byte{0xff}, []byte{0x0000000000000000000000000000000000000000}, []byte{0x0000000000000000000000000000000000000000000000000000000000000000}, crypto.Keccak256([]byte{0x00}))[12:])
		fmt.Println("Keccak256(0x00): ", crypto.Keccak256([]byte{0x00}))
		fmt.Println("Keccak256(0x00): ", string(crypto.Keccak256([]byte{0x00})))
		fmt.Println("Keccak256(0x00): ", hexutil.Encode(crypto.Keccak256([]byte{0x00})))
		fmt.Println("Keccak256(0x0000): ", crypto.Keccak256([]byte{0x000000000000}))
		fmt.Println("Resulting Hash: ", hexutil.Encode(hash.Bytes()))

		// hash := crypto.Keccak256Hash(code)

		// hash2 := crypto.Keccak256([]byte{0xff}, caller.Address().Bytes(), common.BigToAddress(salt).Bytes(), crypto.Keccak256Hash(code).Bytes())[12:]
		// addr := common.BytesToAddress(crypto.Keccak256([]byte{0xff}, caller.Address().Bytes(), common.BigToHash(salt).Bytes(), crypto.Keccak256Hash(code).Bytes())[12:])
		addr := common.BytesToAddress(crypto.Keccak256([]byte{0xff}, []byte{0x0000000000000000000000000000000000000000}, []byte{0x0000000000000000000000000000000000000000000000000000000000000000}, crypto.Keccak256([]byte{0x00}))[12:])
		fmt.Println("ADDR: ", hexutil.Encode(addr.Bytes()))
		// crypto.Keccak256Hash()
		// fmt.Println("HASH1: " + crypto.Keccak256Hash(code).Str())
		// fmt.Println("HASH2: " + common.BytesToAddress(hash2).Str())

		// hash := common.BytesToHash((crypto.Keccak256([]byte{0xff}, salt.Bytes(), crypto.Keccak256Hash(code).Bytes())[12:]))
		// addr := common.BytesToAddress((crypto.Keccak256([]byte{0xff}, salt.Bytes(), crypto.Keccak256Hash(code).Bytes())[12:]))
		// ret, address, err = exec(env, caller, &addr, nil, hash, nil, code, gas, gasPrice, value, false)
		ret, address, err = create(env, caller, &addr, nil, hash, nil, code, gas, gasPrice, value, false)
		fmt.Println("ADDRESS: " + hexutil.Encode(address.Bytes()))*/
>>>>>>> Implemented create function for Create2 contracts
	// Here we get an error if we run into maximum stack depth,
	// See: https://github.com/ethereum/yellowpaper/pull/131
	// and YP definitions for CREATE

	//if there's an error we return nothing
	if err != nil && err != vm.ErrRevert {
		return nil, address, err
	}
	return ret, address, err
}

// create creates a new contract using code as deployment code.
func create(env vm.Environment, caller vm.ContractRef, address, codeAddr *common.Address, codeHash common.Hash, input, code []byte, gas, gasPrice, value *big.Int, readOnly bool) ([]byte, common.Address, error) {
	evm := env.Vm()
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if env.Depth() > callCreateDepthMax {
		caller.ReturnGas(gas, gasPrice)

		return nil, common.Address{}, errCallCreateDepth
	}
<<<<<<< HEAD
=======
	// if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
	// 	return nil, common.Address{}, gas, ErrInsufficientBalance
	// }
>>>>>>> Implemented create function for Create2 contracts
	if !env.CanTransfer(caller.Address(), value) {
		caller.ReturnGas(gas, gasPrice)

		return nil, common.Address{}, ErrInsufficientFunds
	}
<<<<<<< HEAD
=======
	////////////////////////////////////////////////////////////////////////////////////

>>>>>>> Implemented create function for Create2 contracts
	nonce := env.Db().GetNonce(caller.Address())
	env.Db().SetNonce(caller.Address(), nonce+1)

	// Ensure there's no existing contract already at the designated address
	// contractHash := evm.StateDB.GetCodeHash(address)
	contractHash := env.Db().GetCodeHash(*address)
	if env.Db().GetNonce(*address) != 0 || env.Db().GetCode(*address) != nil || (contractHash != (common.Hash{}) && contractHash != crypto.Keccak256Hash(nil)) {
		return nil, common.Address{}, ErrContractAddressCollision
	}
	// Create a new account on the state
	snapshot := env.SnapshotDatabase()
<<<<<<< HEAD

	//Create account with address
	to := env.Db().CreateAccount(*address)
=======
	// env.Db().CreateAccount(*address)

	//TODO: Wrap with check of EIP
	// if env.RuleSet().IsAtlantis(env.BlockNumber()) {
	// 	env.Db().SetNonce(*address, 1)
	// }

	// if evm.ChainConfig().IsEIP158(evm.BlockNumber) {
	// 	evm.StateDB.SetNonce(address, 1)
	// }
	// evm.Transfer(evm.StateDB, caller.Address(), address, value)

	// var from = env.Db().GetAccount(caller.Address())
	var to vm.Account

	//Create account with address
	to = env.Db().CreateAccount(*address)
>>>>>>> Implemented create function for Create2 contracts

	if env.RuleSet().IsAtlantis(env.BlockNumber()) {
		env.Db().SetNonce(*address, 1)
	}

	env.Transfer(env.Db().GetAccount(caller.Address()), env.Db().GetAccount(*address), value)

<<<<<<< HEAD
=======
	// env.Transfer(caller.)

	// contract := vm.NewContract(caller, to, value, gas, gasPrice)

	// env.Transfer(from, to, value)

>>>>>>> Implemented create function for Create2 contracts
	// Initialise a new contract and set the code that is to be used by the EVM.
	// The contract is a scoped environment for this execution context only.
	contract := vm.NewContract(caller, to, value, gas, gasPrice)
	contract.SetCallCode(codeAddr, codeHash, code)
<<<<<<< HEAD
	defer contract.Finalise()

=======
	fmt.Println("CODE!!:", code)
	defer contract.Finalise()

	// if evm.vmConfig.NoRecursion && evm.depth > 0 {
	// 	return nil, address, gas, nil
	// }

	// if evm.vmConfig.Debug && evm.depth == 0 {
	// 	evm.vmConfig.Tracer.CaptureStart(caller.Address(), address, true, codeAndHash.code, gas, value)
	// }
	// start := time.Now()

>>>>>>> Implemented create function for Create2 contracts
	ret, err := evm.Run(contract, input, readOnly)

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := len(ret) > maxCodeSize && env.RuleSet().IsAtlantis(env.BlockNumber())
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		// createDataGas := big.NewInt(int64(len(ret)) * int64(params.CreateDataGas))
		createDataGas := big.NewInt(int64(len(ret)))
		createDataGas.Mul(createDataGas, big.NewInt(200))
		if contract.UseGas(createDataGas) {
			env.Db().SetCode(*address, ret)
		} else {
			err = vm.CodeStoreOutOfGasError
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && (env.RuleSet().IsHomestead(env.BlockNumber()) || err != ErrCodeStoreOutOfGas)) {
		env.RevertToSnapshot(snapshot)
		if err != vm.ErrRevert {
			contract.UseGas(contract.Gas)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}

	// if evm.vmConfig.Debug && env.Depth() == 0 {
	// 	evm.vmConfig.Tracer.CaptureEnd(ret, gas-contract.Gas, time.Since(start), err)
	// }
	return ret, *address, err

}

func exec(env vm.Environment, caller vm.ContractRef, address, codeAddr *common.Address, codeHash common.Hash, input, code []byte, gas, gasPrice, value *big.Int, readOnly bool) (ret []byte, addr common.Address, err error) {
	evm := env.Vm()
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if env.Depth() > callCreateDepthMax {
		caller.ReturnGas(gas, gasPrice)

		return nil, common.Address{}, errCallCreateDepth
	}

	if !env.CanTransfer(caller.Address(), value) {
		caller.ReturnGas(gas, gasPrice)

		return nil, common.Address{}, ValueTransferErr("insufficient funds to transfer value. Req %v, has %v", value, env.Db().GetBalance(caller.Address()))
	}

	var createAccount bool = false
	if address == nil {
		fmt.Println("Creating new Address")
		// Create a new account on the state
		nonce := env.Db().GetNonce(caller.Address())
		env.Db().SetNonce(caller.Address(), nonce+1)
		addr = crypto.CreateAddress(caller.Address(), nonce)
		address = &addr
		createAccount = true
	}

	snapshotPreTransfer := env.SnapshotDatabase()
	var (
		from = env.Db().GetAccount(caller.Address())
		to   vm.Account
	)

	if createAccount {
		to = env.Db().CreateAccount(*address)

		if env.RuleSet().IsAtlantis(env.BlockNumber()) {
			env.Db().SetNonce(*address, 1)
		}
	} else {
		if !env.Db().Exist(*address) {
			//no account may change state from non-existent to existent-but-empty. Refund sender.
			if vm.PrecompiledAtlantis[(*address).Str()] == nil && env.RuleSet().IsAtlantis(env.BlockNumber()) && value.BitLen() == 0 {
				caller.ReturnGas(gas, gasPrice)
				return nil, common.Address{}, nil
			}
			to = env.Db().CreateAccount(*address)
		} else {
			to = env.Db().GetAccount(*address)
		}
	}

	env.Transfer(from, to, value)

	// initialise a new contract and set the code that is to be used by the
	// EVM. The contract is a scoped environment for this execution context
	// only.
	contract := vm.NewContract(caller, to, value, gas, gasPrice)

	// if codeAddr == nil {
	// 	fmt.Println("CODEADDR IS NIL")
	// }
	contract.SetCallCode(codeAddr, codeHash, code)
	defer contract.Finalise()

	fmt.Println("Contract.Caller: ", hexutil.Encode(contract.Caller().Bytes()))
	// fmt.Println("Contract.to: ", contract.toAddr)
	fmt.Println("Contract.Value: ", contract.Value())
	// fmt.Println("Contract.Gas: ", hexutil.Encode(contract.Gas.Bytes()))
	if address == nil {
		fmt.Println("CODEADDR NIL")
	} else {
		fmt.Println("CODEADDR NOT NIL")
		fmt.Println("ADDRESS: ", hexutil.Encode(address.Bytes()))
	}

	// fmt.Println("Contract.codeAddr: ", *codeAddr)
	// fmt.Println("Contract.Value: ", contract.g)
	ret, err = evm.Run(contract, input, readOnly)
	fmt.Println("ret: ", ret)
	fmt.Println("err: ", err)

	maxCodeSizeExceeded := len(ret) > maxCodeSize && env.RuleSet().IsAtlantis(env.BlockNumber())
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && createAccount && !maxCodeSizeExceeded {
		dataGas := big.NewInt(int64(len(ret)))
		// create data gas
		dataGas.Mul(dataGas, big.NewInt(200))
		if contract.UseGas(dataGas) {
			env.Db().SetCode(*address, ret)
		} else {
			err = vm.CodeStoreOutOfGasError
		}
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if createAccount && maxCodeSizeExceeded || (err != nil && (env.RuleSet().IsHomestead(env.BlockNumber()) || err != vm.CodeStoreOutOfGasError)) {
		env.RevertToSnapshot(snapshotPreTransfer)
		if err != vm.ErrRevert {
			contract.UseGas(contract.Gas)
		}
	}

	// When there are no errors but the maxCodeSize is still exceeded, makes more sense than just failing
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}

	return ret, addr, err
}

func execDelegateCall(env vm.Environment, caller vm.ContractRef, originAddr, toAddr, codeAddr *common.Address, codeHash common.Hash, input, code []byte, gas, gasPrice, value *big.Int) (ret []byte, addr common.Address, err error) {
	evm := env.Vm()
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if env.Depth() > callCreateDepthMax {
		caller.ReturnGas(gas, gasPrice)
		return nil, common.Address{}, errCallCreateDepth
	}

	snapshot := env.SnapshotDatabase()

	var to vm.Account
	if !env.Db().Exist(*toAddr) {
		to = env.Db().CreateAccount(*toAddr)
	} else {
		to = env.Db().GetAccount(*toAddr)
	}

	// Iinitialise a new contract and make initialise the delegate values
	contract := vm.NewContract(caller, to, value, gas, gasPrice).AsDelegate()
	contract.SetCallCode(codeAddr, codeHash, code)
	defer contract.Finalise()

	ret, err = evm.Run(contract, input, false)

	if err != nil {
		env.RevertToSnapshot(snapshot)
		if err != vm.ErrRevert {
			contract.UseGas(contract.Gas)
		}
	}

	return ret, addr, err
}

// generic transfer method
func Transfer(from, to vm.Account, amount *big.Int) {
	from.SubBalance(amount)
	to.AddBalance(amount)
}
