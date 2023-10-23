package services

import (
	"context"
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	goethclient "github.com/ethereum/go-ethereum/ethclient"

	"eth-blockchain-service/internal/ethclient"
	"eth-blockchain-service/internal/models"
)

type EthClientService interface {
	GetBlock(ctx context.Context, blockNum int) (*models.Block, []*models.Transaction, error)
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
