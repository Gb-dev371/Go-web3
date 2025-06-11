package contract

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Struct correspondente ao retorno da função byAddress
type Lp struct {
	Lp             common.Address
	Symbol         string
	Decimals       uint8
	Liquidity      *big.Int
	Type           int32
	Tick           int32
	SqrtRatio      *big.Int
	Token0         common.Address
	Reserve0       *big.Int
	Staked0        *big.Int
	Token1         common.Address
	Reserve1       *big.Int
	Staked1        *big.Int
	Gauge          common.Address
	GaugeLiquidity *big.Int
	GaugeAlive     bool
	Fee            common.Address
	Bribe          common.Address
	Factory        common.Address
	Emissions      *big.Int
	EmissionsToken common.Address
	PoolFee        *big.Int
	UnstakedFee    *big.Int
	Token0Fees     *big.Int
	Token1Fees     *big.Int
	Nfpm           common.Address
	Alm            common.Address
	Root           common.Address
}

type LpSugar struct {
	Address common.Address
	ABI     abi.ABI
	Client  *ethclient.Client
}

func NewLpSugar(address string, abiPath string, client *ethclient.Client) *LpSugar {
	abiBytes, err := os.ReadFile(abiPath)
	if err != nil {
		log.Fatalf("Erro ao ler ABI: %v", err)
	}

	parsed, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		log.Fatalf("Erro ao parsear ABI: %v", err)
	}

	return &LpSugar{
		Address: common.HexToAddress(address),
		ABI:     parsed,
		Client:  client,
	}
}

// Usando blockNumber como *big.Int nilável
func (p *LpSugar) ByAddress(pool common.Address, blockNumber *big.Int) (*Lp, error) {
	data, err := p.ABI.Pack("byAddress", pool)
	if err != nil {
		return nil, fmt.Errorf("erro ao empacotar byAddress: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &p.Address,
		Data: data,
		Gas:  5_000_000,
	}

	// Aqui está o ajuste: usamos o estado mais recente com blockNumber = nil
	res, err := p.Client.CallContract(context.Background(), msg, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar byAddress: %w", err)
	}

	var out Lp
	err = p.ABI.UnpackIntoInterface(&out, "byAddress", res)
	if err != nil {
		return nil, fmt.Errorf("erro ao desempacotar byAddress: %w", err)
	}

	return &out, nil
}
