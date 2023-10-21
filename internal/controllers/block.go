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
	blockSrv services.BlockService
}

func NewBlockController() (BlockController, error) {
	srv, err := services.GetService()
	if err != nil {
		return nil, err
	}
	return &blockController{blockSrv: srv.Block}, nil
}

func (c *blockController) GetBlocks(ctx *gin.Context) {
	limitN := ctx.Query("n")
	respond(ctx, nil, map[string]string{"n": limitN}, http.StatusOK)
}

func (c *blockController) GetSingleBlock(ctx *gin.Context) {
	id := ctx.Param("id")
	respond(ctx, nil, map[string]string{"id": id}, http.StatusOK)
}
