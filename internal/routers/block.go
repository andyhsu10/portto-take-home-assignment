package routers

import (
	"github.com/gin-gonic/gin"

	"eth-blockchain-service/internal/controllers"
)

func InitBlockRouter(engine *gin.Engine, path string) error {
	ctl, err := controllers.GetController()
	if err != nil {
		return nil
	}
	group := engine.Group(path)

	group.GET("", ctl.Block.GetBlocks)
	group.GET("/:id", ctl.Block.GetSingleBlock)
	return nil
}
