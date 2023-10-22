package app

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"

	"eth-blockchain-service/internal/configs"
)

func Run() {
	conf := configs.GetConfig()
	cl, err := ethclient.Dial(conf.RpcList[0])
	fmt.Println("ethclient", cl, err)

	ctx := context.Background()
	num, err := cl.BlockNumber(ctx)
	fmt.Println("BlockNumber", num, err)

	block, err := cl.BlockByNumber(ctx, big.NewInt(int64(num)))
	fmt.Println("Block", block, err)
	hash, time, parentHash, transactions := block.Hash(), block.Time(), block.ParentHash(), block.Transactions()
	fmt.Println("BlockInfo", hash, time, parentHash, transactions)
}
