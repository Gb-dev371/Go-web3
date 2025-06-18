package utils

import (
	"fmt"
	"math/big"
)

// CalculateAPR calcula o APR baseado nas emissões em USD e no TVL (Total Value Locked).
func CalculateAPR(emissionsUSD, tvlUSD *big.Float) (*big.Float, error) {
	if tvlUSD.Cmp(big.NewFloat(0)) == 0 {
		return nil, fmt.Errorf("TVL não pode ser zero")
	}

	weeks := big.NewFloat(365.0 / 7.0)
	emissionsAnnual := new(big.Float).Mul(emissionsUSD, weeks)
	apr := new(big.Float).Quo(emissionsAnnual, tvlUSD)
	apr.Mul(apr, big.NewFloat(100))

	return apr, nil
}
