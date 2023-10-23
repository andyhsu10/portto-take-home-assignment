package services

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"eth-blockchain-service/internal/databases"
	"eth-blockchain-service/internal/models"
)

type TxnService interface {
	BatchCreateTxns(ctx context.Context, txns []models.Transaction) (*[]models.Transaction, error)
}

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

func (srv *txnService) BatchCreateTxns(ctx context.Context, txns []models.Transaction) (*[]models.Transaction, error) {
	if len(txns) == 0 {
		return &txns, nil
	}

	res := srv.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&txns)
	if res.Error != nil {
		return nil, res.Error
	}

	return &txns, nil
}
