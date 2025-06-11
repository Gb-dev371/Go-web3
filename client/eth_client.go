package client

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func Connect(rpcURL string) *ethclient.Client {
	RpcClient, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Erro ao conectar ao node RPC: %v", err)
	}
	return RpcClient
}
