package services

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	goethclient "github.com/ethereum/go-ethereum/ethclient"

	"eth-blockchain-service/internal/ethclient"
	"eth-blockchain-service/internal/models"
)

type EthClientService interface {
	GetBlock(ctx context.Context, blockNum int) (*models.Block, []*models.Transaction, error)
	GetTxnLogs(ctx context.Context, hash string) ([]*TxnReceiptLog, error)
}

type ethClientService struct {
	client *goethclient.Client
}

func NewEthClientService() (EthClientService, error) {
	client, err := ethclient.GetClient()
	if err != nil {
		return nil, err
	}

	return &ethClientService{client: client}, nil
}

func (srv *ethClientService) GetBlock(ctx context.Context, blockNum int) (*models.Block, []*models.Transaction, error) {
	b, err := srv.client.BlockByNumber(ctx, big.NewInt(int64(blockNum)))
	if err != nil {
		return nil, nil, err
	}

	block := &models.Block{
		Number:     blockNum,
		Hash:       b.Hash().String(),
		Time:       int(b.Time()),
		ParentHash: b.ParentHash().String(),
	}

	transactions := make([]*models.Transaction, len(b.Transactions()))
	for i, t := range b.Transactions() {
		from, _ := types.Sender(types.LatestSignerForChainID(t.ChainId()), t)
		data := hex.EncodeToString(t.Data())
		if len(data) != 0 {
			data = "0x" + data
		}

		transactions[i] = &models.Transaction{
			Hash:        t.Hash().String(),
			From:        from.String(),
			To:          t.To().String(),
			Nonce:       int(t.Nonce()),
			Data:        data,
			Value:       t.Value().String(),
			BlockNumber: blockNum,
		}
	}

	return block, transactions, nil
}

func (srv *ethClientService) GetTxnLogs(ctx context.Context, hash string) ([]*TxnReceiptLog, error) {
	h := common.HexToHash(hash)
	receipt, err := srv.client.TransactionReceipt(ctx, h)
	if err != nil {
		return nil, err
	}

	logs := make([]*TxnReceiptLog, len(receipt.Logs))
	for i, l := range receipt.Logs {
		logs[i] = &TxnReceiptLog{
			Index: int(l.Index),
			Data:  "0x" + hex.EncodeToString(l.Data),
		}
	}

	return logs, nil
}

type TxnReceiptLog struct {
	Index int    `json:"index"`
	Data  string `json:"data"`
}
