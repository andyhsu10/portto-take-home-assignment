package routers

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"eth-blockchain-service/internal/middlewares"
)

var (
	routerInstance *gin.Engine
)

func GetRouter() (*gin.Engine, error) {
	if routerInstance == nil {
		instance, err := newRouter()
		if err != nil {
			return nil, err
		}
		routerInstance = instance
	}
	return routerInstance, nil
}

func newRouter() (*gin.Engine, error) {
	engine := gin.Default()
	middleware, err := middlewares.GetMiddleware()
	if err != nil {
		return nil, err
	}

	engine.Use(middleware.Cors.Cors())

	engine.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"data": "Hello World!",
		})
	})
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = InitBlockRouter(engine, "blocks")
	if err != nil {
		return nil, err
	}

	err = InitTxnRouter(engine, "transaction")
	if err != nil {
		return nil, err
	}

	return engine, nil
}
