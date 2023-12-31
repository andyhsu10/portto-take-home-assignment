package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func respond(ctx *gin.Context, err error, data interface{}, httpCode int) {
	if err != nil {
		log.Print(err)
		ctx.JSON(httpCode, map[string]interface{}{"message": err.Error(), "data": nil})
		return
	}

	if data != nil {
		ctx.JSON(httpCode, data)
	} else {
		ctx.Status(httpCode)
	}

}
