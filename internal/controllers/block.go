package controllers

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"eth-blockchain-service/internal/configs"
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
	n, err := strconv.Atoi(ctx.Query("n"))
	if err != nil {
		respond(ctx, nil, nil, http.StatusBadRequest)
		return
	}

	if n > configs.GetConfig().MaxN {
		n = configs.GetConfig().MaxN
	}

	blocks, err := c.blockSrv.GetLatestNBlocks(ctx, n)
	if err != nil {
		respond(ctx, nil, nil, http.StatusInternalServerError)
		return
	}

	respond(ctx, nil, blocks, http.StatusOK)
}

func (c *blockController) GetSingleBlock(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respond(ctx, nil, nil, http.StatusBadRequest)
		return
	}

	block, err := c.blockSrv.GetSingleBlock(context.Background(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			respond(ctx, nil, nil, http.StatusNotFound)
			return
		}
		respond(ctx, nil, nil, http.StatusInternalServerError)
		return
	}

	respond(ctx, nil, block, http.StatusOK)
}
