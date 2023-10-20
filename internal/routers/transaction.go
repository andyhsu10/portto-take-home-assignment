package routers

import (
	"github.com/gin-gonic/gin"

	"eth-blockchain-service/internal/controllers"
)

func InitTxnRouter(engine *gin.Engine, path string) error {
	ctl, err := controllers.GetController()
	if err != nil {
		return nil
	}
	group := engine.Group(path)

	group.GET("/:txHash", ctl.Txn.GetSingleTransaction)
	return nil
}
