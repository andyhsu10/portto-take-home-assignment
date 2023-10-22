package ethclient

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"

	"eth-blockchain-service/internal/configs"
)

// set up a singleton instance of the database
var clientInstance *ethclient.Client

func GetClient() (instance *ethclient.Client, err error) {
	if clientInstance == nil {
		instance, err = newClient()
		if err != nil {
			return nil, err
		}
		clientInstance = instance
	}
	return clientInstance, nil
}

func newClient() (*ethclient.Client, error) {
	RpcList := configs.GetConfig().RpcList

	for _, rpc := range RpcList {
		cl, err := ethclient.Dial(rpc)
		if err == nil {
			return cl, nil
		}
	}

	return nil, errors.New("failed to connect to a valid RPC")
}
