package utils

import (
	"fmt"
	"math"
	"math/big"

	"go_web3/contract"
)

// Constants
var Q96 = new(big.Int).Exp(big.NewInt(2), big.NewInt(96), nil) // 2^96

// SqrtPriceX96ToPrice converte sqrt_price_x96 para o preço normalizado entre os tokens
func SqrtPriceX96ToPrice(sqrtPriceX96 *big.Int, decimals0, decimals1 int) *big.Float {
	// Convert sqrtPriceX96 para *big.Float
	sqrtPrice := new(big.Float).SetInt(sqrtPriceX96)

	// Q96 = 2^96
	q96 := new(big.Float).SetFloat64(math.Pow(2, 96))

	// sqrtPriceNormalized = sqrtPriceX96 / Q96
	sqrtPriceNormalized := new(big.Float).Quo(sqrtPrice, q96)

	// price = (sqrtPriceNormalized)^2
	price := new(big.Float).Mul(sqrtPriceNormalized, sqrtPriceNormalized)

	// Ajuste de decimais: * 10^(decimals0 - decimals1)
	decimalAdjustment := new(big.Float).SetFloat64(math.Pow10(decimals0 - decimals1))
	priceAdjusted := new(big.Float).Mul(price, decimalAdjustment)

	return priceAdjusted
}

// Converte um tick para sqrtPriceX96 como big.Float
func TickToSqrtPriceX96(tick int) *big.Float {
	base := 1.0001
	exponent := float64(tick)
	price := math.Pow(base, exponent)
	sqrtPrice := math.Sqrt(price)
	sqrtPriceX96 := new(big.Float).Mul(big.NewFloat(sqrtPrice), new(big.Float).SetInt(Q96))
	return sqrtPriceX96
}

// GetAmount calcula amount0, amount1 e se está dentro do intervalo
func GetAmount(
	tickLower, tickUpper int,
	liquidity *big.Int,
	sqrtPriceX96 *big.Int,
) (*big.Int, *big.Int) {

	// Q96 = 2^96
	Q96 := math.Pow(2, 96)
	sqrtPX96Float := new(big.Float).SetInt(sqrtPriceX96)
	sqrtPFloat := new(big.Float).Quo(sqrtPX96Float, big.NewFloat(Q96))
	sqrtP64, _ := sqrtPFloat.Float64()

	// current_tick = log((sqrtP^2)) / log(1.0001)
	currentTick := int(math.Log(sqrtP64*sqrtP64) / math.Log(1.0001))

	sqrtRatioL := math.Sqrt(math.Pow(1.0001, float64(tickLower)))
	sqrtRatioU := math.Sqrt(math.Pow(1.0001, float64(tickUpper)))
	sqrtP := sqrtP64

	var amount0, amount1 *big.Int

	liqFloat := new(big.Float).SetInt(liquidity)

	switch {
	case currentTick < tickLower:
		numer := sqrtRatioU - sqrtRatioL
		denom := sqrtRatioL * sqrtRatioU
		res := new(big.Float).Quo(
			new(big.Float).Mul(liqFloat, big.NewFloat(numer)),
			big.NewFloat(denom),
		)
		amount0 = new(big.Int)
		res.Int(amount0)
		amount1 = big.NewInt(0)

	case currentTick >= tickUpper:
		diff := sqrtRatioU - sqrtRatioL
		res := new(big.Float).Mul(liqFloat, big.NewFloat(diff))
		amount1 = new(big.Int)
		res.Int(amount1)
		amount0 = big.NewInt(0)

	case currentTick >= tickLower && currentTick < tickUpper:
		diff0 := sqrtRatioU - sqrtP
		den0 := sqrtP * sqrtRatioU
		val0 := new(big.Float).Quo(
			new(big.Float).Mul(liqFloat, big.NewFloat(diff0)),
			big.NewFloat(den0),
		)
		amount0 = new(big.Int)
		val0.Int(amount0)

		diff1 := sqrtP - sqrtRatioL
		val1 := new(big.Float).Mul(liqFloat, big.NewFloat(diff1))
		amount1 = new(big.Int)
		val1.Int(amount1)

	}

	return amount0, amount1
}

func GetCurrentAmountsInLiquidity(
	p *contract.CLPool,
	liquidity *big.Int,
	blockNumber *big.Int,
	decimals0 int,
	decimals1 int,
) (float64, float64, error) {

	tickSpacingBig, err := p.TickSpacing(blockNumber)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao obter tickSpacing: %v", err)
	}

	spacing := int(tickSpacingBig.Int64())

	slot0, err := p.Slot0(blockNumber)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao obter slot0: %v", err)
	}
	sqrtPriceX96 := slot0.SqrtPriceX96
	currentTick := int(slot0.Tick) // já é int32, pode converter para int

	decrement := currentTick % spacing
	increment := spacing - decrement
	tickLower := currentTick - decrement
	tickUpper := currentTick + increment

	amount0, amount1 := GetAmount(tickLower, tickUpper, liquidity, sqrtPriceX96)

	amount0Float := new(big.Float).Quo(new(big.Float).SetInt(amount0), big.NewFloat(math.Pow10(decimals0)))
	amount1Float := new(big.Float).Quo(new(big.Float).SetInt(amount1), big.NewFloat(math.Pow10(decimals1)))

	a0, _ := amount0Float.Float64()
	a1, _ := amount1Float.Float64()

	return a0, a1, nil
}
