package services

import (
	"gorm.io/gorm"

	"eth-blockchain-service/internal/databases"
)

type BlockService interface{}

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
