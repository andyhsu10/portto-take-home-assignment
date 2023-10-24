package services

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"eth-blockchain-service/internal/databases"
	"eth-blockchain-service/internal/models"
)

type BlockService interface {
	GetSingleBlock(ctx context.Context, blockNum int) (*SingleBlockResponse, error)
	GetLatestNBlocks(ctx context.Context, num int) (*BlocksResponse, error)
	GetBlockNumbers(ctx context.Context, nums *[]int) (*[]int, error)
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

func (srv *blockService) GetSingleBlock(ctx context.Context, blockNum int) (*SingleBlockResponse, error) {
	block := &models.Block{}
	res := srv.db.
		Where("number = ?", blockNum).
		Preload("Transactions", func(db *gorm.DB) *gorm.DB {
			return db.Select("Hash", "BlockNumber")
		}).
		First(&block)

	if res.Error != nil {
		return nil, res.Error
	}

	txids := make([]string, len(block.Transactions))
	for i, t := range block.Transactions {
		txids[i] = t.Hash
	}

	response := &SingleBlockResponse{
		Block:        *block,
		Transactions: txids,
	}
	return response, nil
}

func (srv *blockService) GetLatestNBlocks(ctx context.Context, num int) (*BlocksResponse, error) {
	blocks := make([]models.Block, 0)
	res := srv.db.Limit(num).Order("number DESC").Find(&blocks)

	if res.Error != nil {
		return nil, res.Error
	}

	response := &BlocksResponse{
		Blocks: blocks,
	}

	return response, nil
}

func (srv *blockService) GetBlockNumbers(ctx context.Context, nums *[]int) (*[]int, error) {
	blocks := make([]*models.Block, 0)
	res := srv.db.Where("number IN ?", *nums).Select("Number").Find(&blocks)
	if res.Error != nil {
		return nil, res.Error
	}

	result := make([]int, len(blocks))
	for i, b := range blocks {
		result[i] = int(b.Number)
	}

	return &result, nil
}

func (srv *blockService) CreateBlock(ctx context.Context, block models.Block) (*models.Block, error) {
	res := srv.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&block)
	if res.Error != nil {
		return nil, res.Error
	}

	return &block, nil
}

type BlocksResponse struct {
	Blocks []models.Block `json:"blocks"`
}

type SingleBlockResponse struct {
	models.Block
	Transactions []string `json:"transactions"`
}
