package app

import (
	"context"
	"fmt"
	"math/big"

	"eth-blockchain-service/internal/ethclient"
)

func Run() {
	cl, err := ethclient.GetClient()
	fmt.Println("ethclient", cl, err)

	ctx := context.Background()
	num, err := cl.BlockNumber(ctx)
	fmt.Println("BlockNumber", num, err)

	block, err := cl.BlockByNumber(ctx, big.NewInt(int64(num)))
	fmt.Println("Block", block, err)
	hash, time, parentHash, transactions := block.Hash(), block.Time(), block.ParentHash(), block.Transactions()
	fmt.Println("BlockInfo", hash, time, parentHash, transactions)
}
