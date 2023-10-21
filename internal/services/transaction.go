package services

import (
	"gorm.io/gorm"

	"eth-blockchain-service/internal/databases"
)

type TxnService interface{}

type txnService struct {
	db *gorm.DB
}

func NewTxnService() (TxnService, error) {
	db, err := databases.GetDB()
	if err != nil {
		return nil, err
	}

	return &txnService{db: db}, nil
}
