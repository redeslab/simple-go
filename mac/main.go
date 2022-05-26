package main

import "C"
import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/redeslab/go-simple/account"
)

var (
	abiUintType, _   = abi.NewType("uint256", "", nil)
	abiAddrType, _   = abi.NewType("address", "", nil)
	abiStrType, _    = abi.NewType("string", "", nil)
	abiByte32Type, _ = abi.NewType("bytes32", "", nil)
	abiMicroTxArgs   = abi.Arguments{
		{Type: abiAddrType},
		{Type: abiAddrType},
		{Type: abiAddrType},
		{Type: abiAddrType},
		{Type: abiUintType},
		{Type: abiUintType},
		{Type: abiUintType},
		{Type: abiUintType},
	}
	abiPrefixHashArgs = abi.Arguments{
		{Type: abiStrType},
		{Type: abiByte32Type},
	}
)

func main() {
	pub1 := account.ID("HO6Qh2sabkHhrh7UqmAUzZ3Epe3JgcHusLB8QVg54WGSYb").ToArray()
	fmt.Println(pub1)
}
