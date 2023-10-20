package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eth-blockchain-service/internal/services"
)

type TxnController interface {
	GetSingleTransaction(ctx *gin.Context)
}

type txnController struct {
	txn services.TxnService
}

func NewTxnController() (TxnController, error) {
	srv, err := services.GetService()
	if err != nil {
		return nil, err
	}

	return &txnController{txn: srv.Txn}, nil
}

func (controller *txnController) GetSingleTransaction(ctx *gin.Context) {
	respond(ctx, nil, map[string]string{"controller": "GetSingleTransaction"}, http.StatusOK)
}
