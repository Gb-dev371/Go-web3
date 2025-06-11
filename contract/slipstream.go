package contract

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type CLPool struct {
	Address common.Address
	ABI     abi.ABI
	Client  *ethclient.Client
}

func NewCLPool(address string, abiPath string, client *ethclient.Client) *CLPool {
	// Carrega ABI
	abiBytes, err := os.ReadFile(abiPath)
	if err != nil {
		log.Fatalf("Erro ao ler ABI: %v", err)
	}
	var parsed abi.ABI
	err = json.Unmarshal(abiBytes, &parsed)
	if err != nil {
		parsed, err = abi.JSON(strings.NewReader(string(abiBytes)))
		if err != nil {
			log.Fatalf("Erro ao parsear ABI: %v", err)
		}
	}

	return &CLPool{
		Address: common.HexToAddress(address),
		ABI:     parsed,
		Client:  client,
	}
}

func (p *CLPool) callUint256Method(method string, blockNumber *big.Int) *big.Int {
	data, err := p.ABI.Pack(method)
	if err != nil {
		log.Fatalf("Erro ao empacotar método %s: %v", method, err)
	}

	msg := ethereum.CallMsg{
		To:   &p.Address,
		Data: data,
	}

	res, err := p.Client.CallContract(context.Background(), msg, blockNumber)
	if err != nil {
		log.Fatalf("Erro ao chamar método %s: %v", method, err)
	}

	var output *big.Int
	err = p.ABI.UnpackIntoInterface(&output, method, res)
	if err != nil {
		log.Fatalf("Erro ao decodificar retorno %s: %v", method, err)
	}

	return output
}

// Exemplo de métodos públicos

func (p *CLPool) StakedLiquidity(blockNumber *big.Int) *big.Int {
	return p.callUint256Method("stakedLiquidity", blockNumber)
}

func (p *CLPool) Fee(blockNumber *big.Int) *big.Int {
	return p.callUint256Method("fee", blockNumber)
}

func (p *CLPool) RewardRate(blockNumber *big.Int) *big.Int {
	return p.callUint256Method("rewardRate", blockNumber)
}
