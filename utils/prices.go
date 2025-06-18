package utils

import (
	"fmt"
	"log"
	"math/big"

	"go_web3/contract"

	"github.com/ethereum/go-ethereum/ethclient"
)

func GetAeroPriceInUSDC(rpcURL string, blockNumber *big.Int) (*big.Float, error) {
	// Endereço da pool AERO/USDC
	poolAddress := "0xBE00fF35AF70E8415D0eB605a286D8A45466A4c1"
	abiPath := "abis/slipstream_pool_abi.json"

	// Conexão com o cliente
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao RPC: %v", err)
	}

	// Instanciar a pool
	pool := contract.NewCLPool(poolAddress, abiPath, client)

	// Obter dados do slot0
	slot0, err := pool.Slot0(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter slot0 da pool: %v", err)
	}

	// Calcular preço do token0 em termos do token1
	priceToken0 := SqrtPriceX96ToPrice(slot0.SqrtPriceX96, 6, 18)

	// Inverter para obter o preço do token1 (AERO) em token0 (USDC)
	priceToken1 := new(big.Float).Quo(big.NewFloat(1), priceToken0)

	return priceToken1, nil
}

func ExampleUsage() {
	rpcURL := "https://mainnet.base.org"

	price, err := GetAeroPriceInUSDC(rpcURL, nil)
	if err != nil {
		log.Fatalf("Erro ao obter preço do AERO: %v", err)
	}

	fmt.Printf("Preço do AERO em USDC: %f\n", price)
}
