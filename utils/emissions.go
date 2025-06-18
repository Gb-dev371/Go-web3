package utils

import (
	"fmt"
	"log"
	"math/big"

	"go_web3/contract"

	"github.com/ethereum/go-ethereum/ethclient"
)

// GetWeeklyEmissionsUSD calcula as emissões semanais em USD de uma pool CL usando o rewardRate e o preço do AERO
func GetWeeklyEmissionsUSD(
	poolAddress string,
	abiPath string,
	rpcURL string,
	decimalsRewardToken int,
) (*big.Float, error) {

	// Conectar RPC
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao RPC: %v", err)
	}

	// Instanciar contrato da pool
	pool := contract.NewCLPool(poolAddress, abiPath, client)

	block := (*big.Int)(nil)

	// Obter rewardRate no bloco
	rewardRate, err := pool.RewardRate(block)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter rewardRate: %v", err)
	}

	// rewardRate (*big.Int) → *big.Float
	rateFloat := new(big.Float).SetInt(rewardRate)

	// Dividir por 10^decimals
	divisor := new(big.Float).SetFloat64(1e18) // padrão ERC20
	ratePerSecond := new(big.Float).Quo(rateFloat, divisor)

	// Multiplicar pelo número de segundos da semana
	secondsInWeek := new(big.Float).SetFloat64(604800)
	emissionsPerWeek := new(big.Float).Mul(ratePerSecond, secondsInWeek)

	fmt.Printf("Emissions por semana: %.6f\n", emissionsPerWeek)

	// Obter preço do AERO em USDC
	priceAERO, err := GetAeroPriceInUSDC(rpcURL, block)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter preço do AERO: %v", err)
	}

	// Calcular emissões semanais em USD
	emissionsUSD := new(big.Float).Mul(emissionsPerWeek, priceAERO)

	return emissionsUSD, nil
}

func GetEmissionsUsd(poolAddress string) *big.Float {
	rpcURL := "https://base.llamarpc.com/sk_llama_983b3eb1da9b648371f7139f9c7f2b63"

	abiPath := "abis/slipstream_pool_abi.json"

	emissionsUSD, err := GetWeeklyEmissionsUSD(poolAddress, abiPath, rpcURL, 18)
	if err != nil {
		log.Fatalf("Erro: %v", err)
	}

	return emissionsUSD
}
