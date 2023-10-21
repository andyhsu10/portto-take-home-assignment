package services

import (
	"context"

	"gorm.io/gorm"

	"eth-blockchain-service/internal/databases"
	"eth-blockchain-service/internal/models"
)

type BlockService interface {
	GetSingleBlock(ctx context.Context, blockNum int) (*models.Block, error)
}

type blockService struct {
	db *gorm.DB
}

func NewBlockService() (BlockService, error) {
	db, err := databases.GetDB()
	if err != nil {
		return nil, err
	}

	return &blockService{db: db}, nil
}

func (srv *blockService) GetSingleBlock(ctx context.Context, blockNum int) (*models.Block, error) {
	block := &models.Block{}
	res := srv.db.
		Where("number = ?", blockNum).
		Preload("Transactions").
		First(&block)

	if res.Error != nil {
		return nil, res.Error
	}
	return block, nil
}
