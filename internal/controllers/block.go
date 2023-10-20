package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"eth-blockchain-service/internal/services"
)

type BlockController interface {
	GetBlocks(ctx *gin.Context)
	GetSingleBlock(ctx *gin.Context)
}

type blockController struct {
	block services.BlockService
}

func NewBlockController() (BlockController, error) {
	srv, err := services.GetService()
	if err != nil {
		return nil, err
	}
	return &blockController{block: srv.Block}, nil
}

func (controller *blockController) GetBlocks(ctx *gin.Context) {
	respond(ctx, nil, map[string]string{"controller": "GetBlocks"}, http.StatusOK)
}

func (controller *blockController) GetSingleBlock(ctx *gin.Context) {
	respond(ctx, nil, map[string]string{"controller": "GetSingleBlock"}, http.StatusOK)
}
