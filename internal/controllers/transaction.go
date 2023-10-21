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
	txnSrv services.TxnService
}

func NewTxnController() (TxnController, error) {
	srv, err := services.GetService()
	if err != nil {
		return nil, err
	}

	return &txnController{txnSrv: srv.Txn}, nil
}

func (c *txnController) GetSingleTransaction(ctx *gin.Context) {
	txHash := ctx.Param("txHash")
	respond(ctx, nil, map[string]string{"txHash": txHash}, http.StatusOK)
}
