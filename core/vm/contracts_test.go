package vm

import (
	"testing"

	"github.com/eth-classic/go-ethereum/common"
)

type precompiledTest struct {
	input, expected string
	gas             uint64
}

func testEcRecover(test precompiledTest, t *testing.T) {
	in := common.Hex2Bytes(test.input)
	if res, err := ecrecoverFunc(in); err != nil {
		t.Error(err)
	} else if common.Bytes2Hex(res) != test.expected {
		t.Errorf("Expected %v, got %v", test.expected, common.Bytes2Hex(res))
	}
}

func TestPrecompiledEcRecover(t *testing.T) {
	test := precompiledTest{
		input:    "38d18acb67d25c8bb9942764b62f18e17054f66a817bd4295423adf9ed98873e000000000000000000000000000000000000000000000000000000000000001b38d18acb67d25c8bb9942764b62f18e17054f66a817bd4295423adf9ed98873e789d1dd423d25f0772d2748d60f7e4b81bb14d086eba8e8e8efb6dcff8a4ae02",
		expected: "000000000000000000000000ceaccac640adf55b2028469bd36ba501f28b699d",
	}

	testEcRecover(test, t)
}
