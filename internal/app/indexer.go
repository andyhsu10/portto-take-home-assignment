package app

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"eth-blockchain-service/internal/ethclient"
)

type ReceiptLog struct {
	Index int    `json:"index"`
	Data  string `json:"data"`
}

func Run() {
	cl, err := ethclient.GetClient()
	fmt.Println("ethclient", cl, err)

	ctx := context.Background()
	num, err := cl.BlockNumber(ctx)
	fmt.Println("BlockNumber", num, err)

	block, err := cl.BlockByNumber(ctx, big.NewInt(int64(num)))
	fmt.Println("Block", block, err)

	for _, t := range block.Transactions() {
		hash := t.Hash()
		receipt, _ := cl.TransactionReceipt(ctx, hash)
		logs := make([]*ReceiptLog, len(receipt.Logs))
		for i, l := range receipt.Logs {
			logs[i] = &ReceiptLog{
				Index: int(l.Index),
				Data:  "0x" + hex.EncodeToString(l.Data),
			}
		}
		jsonLogs, _ := json.Marshal(logs)

		fmt.Println("TXID:", hash, string(jsonLogs))
	}
}
