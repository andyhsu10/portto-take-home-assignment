package controllers

import (
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

// GetBlocks godoc
// @Summary Get latest N blocks (without txn hash)
// @Tags Block
// @Param n query number true "The number of blocks to get"
// @produce application/json
// @Router /blocks [get]
// @Success 200 {array} models.Block
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

// GetSingleBlock godoc
// @Summary Get block information from a block number (with txn hash)
// @Tags Block
// @Param id path string true "The block number"
// @produce application/json
// @Router /blocks/{id} [get]
// @Success 200 {object} services.BlockResponse
func (c *blockController) GetSingleBlock(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		respond(ctx, nil, nil, http.StatusBadRequest)
		return
	}

	block, err := c.blockSrv.GetSingleBlock(ctx, id)
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
