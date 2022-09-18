package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

var clientEthereum *rpc.Client
var contractBalanceof abi.ABI

const ADDRESS_ZERO = "0x0000000000000000000000000000000000000000"
const poolAddress = "0x0d4a11d5EEaaC28EC3F61d100daF4d40471f1852"
const fromTokenAddress = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
const toTokenAddress = "0xdac17f958d2ee523a2206206994597c13d831ec7"

const balanceOfABIJson = `[{"constant":true,"inputs":[{"name": "","type": "address"}],"name": "balanceOf","outputs":[{"name":"","type": "uint256"}],"payable": false,"stateMutability": "view","type": "function"}]`
const inputAmount = 1e18

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	var resultStrFromToken, resultStrToToken string
	var fromTokenAmount, toTokenAmount *big.Int

	contractBalanceof, err = abi.JSON(strings.NewReader(balanceOfABIJson))
	Check(err)

	poolData, err := contractBalanceof.Pack("balanceOf", common.HexToAddress(poolAddress))

	//Create the new RPC Client that is connect to ethereum mainnet
	clientEthereum, err = rpc.DialHTTP("https://mainnet.infura.io/v3/69461b73caed42a399f9bb9202d63a9c")
	Check(err)

	err = clientEthereum.Call(&resultStrFromToken, "eth_call", map[string]interface{}{
		"from": ADDRESS_ZERO,
		"to":   fromTokenAddress,
		"data": hexutil.Bytes(poolData),
	}, "latest")

	if err == nil {
		resultFrom, err := contractBalanceof.Unpack("balanceOf", hexutil.MustDecode(resultStrFromToken))
		Check(err)
		fromTokenAmount = resultFrom[0].(*big.Int)
		fmt.Printf("The balance of fromToken : %s \n", fromTokenAmount.String())

	} else {
		fmt.Println(err)
	}

	err = clientEthereum.Call(&resultStrToToken, "eth_call", map[string]interface{}{
		"from": ADDRESS_ZERO,
		"to":   toTokenAddress,
		"data": hexutil.Bytes(poolData),
	}, "latest")

	if err == nil {
		resultTo, err := contractBalanceof.Unpack("balanceOf", hexutil.MustDecode(resultStrToToken))
		Check(err)
		toTokenAmount = resultTo[0].(*big.Int)
		fmt.Printf("The balance of toToken : %s \n", toTokenAmount.String())
	} else {
		fmt.Println(err)
	}

	var outputAmount *big.Int

	//Calulate the result by the fomular "const = X * Y"
	fromTokenAmount.Add(fromTokenAmount, big.NewInt(inputAmount))
	toTokenAmount.Mul(toTokenAmount, big.NewInt(inputAmount))
	toTokenAmount.Div(toTokenAmount, fromTokenAmount)

	outputAmount = big.NewInt(toTokenAmount.Int64())
	//Print the output balance of toToken
	fmt.Printf("The result is %s \n", outputAmount.String())

}
