package utils

import (
	"fmt"
	"math/big"

	"go_web3/contract"
)

// GetTVL calculates the USD value of the liquidity in the pool.
func GetTVL(
	p *contract.CLPool,
	liquidity *big.Int,
	blockNumber *big.Int,
	decimals0 int,
	decimals1 int,
) (*big.Float, error) {
	// 1. Calcula amount0 e amount1
	amount0Float, amount1Float, err := GetCurrentAmountsInLiquidity(p, liquidity, blockNumber, decimals0, decimals1)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular os amounts atuais: %v", err)
	}

	// Converte para *big.Float
	amount0 := big.NewFloat(amount0Float)
	amount1 := big.NewFloat(amount1Float)

	// 2. Calcula o preço do token1 em relação ao token0
	slot0, err := p.Slot0(blockNumber)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter slot0 para calcular o preço: %v", err)
	}
	priceToken0 := SqrtPriceX96ToPrice(slot0.SqrtPriceX96, decimals0, decimals1)
	priceToken1 := new(big.Float).Quo(big.NewFloat(1), priceToken0)

	// 3. TVL = amount0 * 1 + amount1 * priceToken1
	tvlToken1 := new(big.Float).Mul(amount1, priceToken1)
	tvlUSD := new(big.Float).Add(amount0, tvlToken1)

	return tvlUSD, nil
}
