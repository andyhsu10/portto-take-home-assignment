package app

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"eth-blockchain-service/internal/ethclient"

	"github.com/ethereum/go-ethereum/core/types"
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
	hash, time, parentHash, transactions := block.Hash(), block.Time(), block.ParentHash(), block.Transactions()
	fmt.Println("BlockInfo", hash, time, parentHash, transactions)

	for _, t := range transactions {
		hash := t.Hash()
		receipt, _ := cl.TransactionReceipt(ctx, hash)
		from, _ := types.Sender(types.LatestSignerForChainID(t.ChainId()), t)
		to := t.To()
		nonce := t.Nonce()
		// FIXME: data might be empty
		data := "0x" + hex.EncodeToString(t.Data())
		value := t.Value()
		logs := make([]*ReceiptLog, len(receipt.Logs))
		for i, l := range receipt.Logs {
			logs[i] = &ReceiptLog{
				Index: int(l.Index),
				Data:  "0x" + hex.EncodeToString(l.Data),
			}
		}
		jsonLogs, _ := json.Marshal(logs)

		fmt.Println("TXID:", hash, from.String(), to.String(), nonce, data, value, string(jsonLogs))
	}
}
