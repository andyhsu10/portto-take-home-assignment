package services

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"eth-blockchain-service/internal/databases"
	"eth-blockchain-service/internal/models"
)

type BlockService interface {
	GetSingleBlock(ctx context.Context, blockNum int) (*models.Block, error)
	GetLatestNBlock(ctx context.Context, num int) (*[]int, error)
	CreateBlock(ctx context.Context, block models.Block) (*models.Block, error)
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

func (srv *blockService) GetLatestNBlock(ctx context.Context, num int) (*[]int, error) {
	blocks := make([]*models.Block, 0)
	res := srv.db.Select("Number").Limit(num).Find(&blocks)
	if res.Error != nil {
		return nil, res.Error
	}

	nums := make([]int, len(blocks))
	for i, b := range blocks {
		nums[i] = int(b.Number)
	}

	return &nums, nil
}

func (srv *blockService) CreateBlock(ctx context.Context, block models.Block) (*models.Block, error) {
	res := srv.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&block)
	if res.Error != nil {
		return nil, res.Error
	}

	return &block, nil
}
