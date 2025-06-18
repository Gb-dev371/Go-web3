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
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Slot0 struct {
	SqrtPriceX96               *big.Int
	Tick                       int32 // Solidity int24 → Go int32
	ObservationIndex           uint16
	ObservationCardinality     uint16
	ObservationCardinalityNext uint16
	FeeProtocol                uint8
}

type CLPool struct {
	Address common.Address
	ABI     abi.ABI
	Client  *ethclient.Client
}

func NewCLPool(address string, abiPath string, client *ethclient.Client) *CLPool {
	abiBytes, err := os.ReadFile(abiPath)
	if err != nil {
		log.Fatalf("Erro ao ler arquivo ABI: %v", err)
	}

	parsed, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		log.Fatalf("Erro ao parsear ABI: %v", err)
	}

	return &CLPool{
		Address: common.HexToAddress(address),
		ABI:     parsed,
		Client:  client,
	}
}

func (p *CLPool) callSingleUint256(method string, blockNumber *big.Int) (*big.Int, error) {
	data, err := p.ABI.Pack(method)
	if err != nil {
		return nil, fmt.Errorf("erro ao empacotar método %s: %w", method, err)
	}

	msg := ethereum.CallMsg{
		To:   &p.Address,
		Data: data,
	}

	res, err := p.Client.CallContract(context.Background(), msg, blockNumber)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar chamada para %s: %w", method, err)
	}

	var output *big.Int
	err = p.ABI.UnpackIntoInterface(&output, method, res)
	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar retorno de %s: %w", method, err)
	}

	return output, nil
}

func (p *CLPool) StakedLiquidity(blockNumber *big.Int) (*big.Int, error) {
	return p.callSingleUint256("stakedLiquidity", blockNumber)
}

func (p *CLPool) Fee(blockNumber *big.Int) (*big.Int, error) {
	return p.callSingleUint256("fee", blockNumber)
}

func (p *CLPool) RewardRate(blockNumber *big.Int) (*big.Int, error) {
	return p.callSingleUint256("rewardRate", blockNumber)
}

func (p *CLPool) TickSpacing(blockNumber *big.Int) (*big.Int, error) {
	return p.callSingleUint256("tickSpacing", blockNumber)
}

func (p *CLPool) Slot0(blockNumber *big.Int) (*Slot0, error) {
	contract := bind.NewBoundContract(p.Address, p.ABI, p.Client, nil, nil)

	var out []interface{}
	err := contract.Call(&bind.CallOpts{
		BlockNumber: blockNumber,
		Context:     context.Background(),
	}, &out, "slot0")

	if err != nil {
		return nil, fmt.Errorf("erro ao chamar slot0: %w", err)
	}

	if len(out) != 6 {
		return nil, fmt.Errorf("esperado 6 valores em slot0, obtido %d", len(out))
	}

	result := &Slot0{
		SqrtPriceX96:               out[0].(*big.Int),
		Tick:                       int32(out[1].(*big.Int).Int64()), // CORRIGIDO AQUI
		ObservationIndex:           out[2].(uint16),
		ObservationCardinality:     out[3].(uint16),
		ObservationCardinalityNext: out[4].(uint16),
	}

	return result, nil
}
