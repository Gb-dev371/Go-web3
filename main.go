package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"go_web3/client"
	"go_web3/contract"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getStakedLiquidity(rpcClient *ethclient.Client, poolAddress string, abiPath string, blockNumber *big.Int) *big.Int {
	clPool := contract.NewCLPool(poolAddress, abiPath, rpcClient)
	return clPool.StakedLiquidity(blockNumber)
}

func getLpSugarData(rpcClient *ethclient.Client, contractAddress, abiPath string, poolAddress common.Address, blockNumber *big.Int) *contract.Lp {
	lpSugar := contract.NewLpSugar(contractAddress, abiPath, rpcClient)
	lpData, err := lpSugar.ByAddress(poolAddress, blockNumber)
	if err != nil {
		log.Fatalf("Erro ao buscar dados da pool via LpSugar: %v", err)
	}
	return lpData
}

func main() {
	rpcURL := "https://base.llamarpc.com"
	rpcClient := client.Connect(rpcURL)
	defer rpcClient.Close()

	// Pegamos o número do bloco apenas para o CLPool
	blockNumber, err := rpcClient.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	blockNum := big.NewInt(int64(blockNumber))

	// === CLPool ===
	poolAddrStr := "0xD43Decd5Df4BDFFd5A4Cf35cA1f9557E33B7246C"
	abiPathCL := "abis/slipstream_pool_abi.json"
	staked := getStakedLiquidity(rpcClient, poolAddrStr, abiPathCL, blockNum)
	fmt.Println("Staked Liquidity:", staked.String())

	// === LpSugar ===
	lpSugarAddress := "0x73ffd28DFde56704F832163e6cD432FCbbD607a1"
	abiPathLp := "abis/lp_sugar_abi.json"
	poolAddress := common.HexToAddress(poolAddrStr)

	// Aqui passamos nil para usar o estado mais recente (igual ao BaseScan)
	lpData := getLpSugarData(rpcClient, lpSugarAddress, abiPathLp, poolAddress, nil)

	fmt.Println("Symbol:", lpData.Symbol)
	fmt.Println("Liquidity:", lpData.Liquidity.String())
	fmt.Println("Gauge:", lpData.Gauge.Hex())
	fmt.Println("Token0:", lpData.Token0.Hex())
	fmt.Println("Reserve0:", lpData.Reserve0.String())
	fmt.Println("Token1:", lpData.Token1.Hex())
	fmt.Println("Reserve1:", lpData.Reserve1.String())
}
