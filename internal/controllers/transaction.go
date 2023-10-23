package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eth-blockchain-service/internal/services"
)

type TxnController interface {
	GetSingleTxn(ctx *gin.Context)
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

func (c *txnController) GetSingleTxn(ctx *gin.Context) {
	txHash := ctx.Param("txHash")
	txn, err := c.txnSrv.GetSingleTxn(ctx, txHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond(ctx, nil, nil, http.StatusNotFound)
			return
		}
		respond(ctx, nil, nil, http.StatusInternalServerError)
		return
	}

	respond(ctx, nil, txn, http.StatusOK)
}
