package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"

	"go_web3/contract"
	"go_web3/utils"
)

func main() {
	// RPC público da Base (substitua se usar outro RPC)
	rpcURL := "https://base.llamarpc.com/sk_llama_983b3eb1da9b648371f7139f9c7f2b63" // ou outro válido

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao nó RPC: %v", err)
	}

	// Endereço da pool (CLPool) e caminho do arquivo ABI
	poolAddress := "0xD43Decd5Df4BDFFd5A4Cf35cA1f9557E33B7246C"
	abiPath := "abis/slipstream_pool_abi.json" // deve conter o método "slot0", "fee", etc.

	pool := contract.NewCLPool(poolAddress, abiPath, client)

	// Bloco mais recente: nil
	block := (*big.Int)(nil)

	fmt.Println("== Informações da pool ==")

	// Staked Liquidity
	stakedLiquidity, err := pool.StakedLiquidity(block)
	if err != nil {
		log.Fatalf("Erro ao buscar stakedLiquidity: %v", err)
	}

	// 4. Chamar função utilitária
	amount0, amount1, err := utils.GetCurrentAmountsInLiquidity(pool, stakedLiquidity, block, 6, 18)
	if err != nil {
		log.Fatalf("Erro ao calcular amounts: %v", err)
	}

	// 5. Imprimir resultado
	fmt.Printf("Amount0: %f\n", amount0)
	fmt.Printf("Amount1: %f\n", amount1)

	tvl, err := utils.GetTVL(pool, stakedLiquidity, block, 6, 18)
	if err != nil {
		fmt.Println("Erro ao calcular TVL:", err)
	} else {
		fmt.Printf("TVL em USD: %.4f\n", tvl)
	}

	// 6. Emissions USD
	emissionsUsd := utils.GetEmissionsUsd("0xD43Decd5Df4BDFFd5A4Cf35cA1f9557E33B7246C")
	fmt.Printf("Emissions usd: %f\n", emissionsUsd)

	apr, err := utils.CalculateAPR(emissionsUsd, tvl)
	if err != nil {
		log.Fatalf("Erro ao calcular APR: %v", err)
	}

	aprFormatted, _ := apr.Float64()
	fmt.Printf("APR da pool: %.2f%%\n", aprFormatted)

}
