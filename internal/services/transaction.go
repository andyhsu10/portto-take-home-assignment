package services

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"eth-blockchain-service/internal/databases"
	"eth-blockchain-service/internal/models"
)

type TxnService interface {
	GetSingleTxn(ctx context.Context, hash string) (*TxnResponse, error)
	BatchCreateTxns(ctx context.Context, txns []models.Transaction) (*[]models.Transaction, error)
	UpdateTxnLogs(ctx context.Context, hash string, logs string) (*models.Transaction, error)
}

type txnService struct {
	db           *gorm.DB
	ethClientSrv EthClientService
}

func NewTxnService() (TxnService, error) {
	db, err := databases.GetDB()
	if err != nil {
		return nil, err
	}

	ethClientSrv, err := NewEthClientService()
	if err != nil {
		return nil, err
	}

	return &txnService{db: db, ethClientSrv: ethClientSrv}, nil
}

func (srv *txnService) GetSingleTxn(ctx context.Context, hash string) (*TxnResponse, error) {
	txn := &models.Transaction{}
	res := srv.db.
		Where("hash = ?", hash).
		First(&txn)

	if res.Error != nil {
		return nil, res.Error
	}

	response := &TxnResponse{
		Transaction: *txn,
	}

	if len(txn.Logs) == 0 {
		logs, err := srv.ethClientSrv.GetTxnLogs(ctx, hash)
		if err == nil {
			response.Logs = *logs
			jsonData, err := json.Marshal(logs)
			if err == nil {
				jsonString := string(jsonData)
				srv.UpdateTxnLogs(ctx, hash, jsonString)
			}
		}
	} else {
		var logs []TxnReceiptLog
		err := json.Unmarshal([]byte(txn.Logs), &logs)
		if err == nil {
			response.Logs = logs
		}
	}

	return response, nil
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

func (srv *txnService) UpdateTxnLogs(ctx context.Context, hash string, logs string) (*models.Transaction, error) {
	txn := &models.Transaction{}
	res := srv.db.Where("hash = ?", hash).First(txn)
	if res.Error != nil {
		return nil, res.Error
	}

	txn.Logs = logs
	// update logs
	res = srv.db.Save(txn)
	if res.Error != nil {
		return nil, res.Error
	}

	return txn, nil
}

type TxnResponse struct {
	models.Transaction
	Logs []TxnReceiptLog `json:"logs"`
}
