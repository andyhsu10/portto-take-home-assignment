package main

import (
	"fmt"
	"log"

	_ "eth-blockchain-service/docs"
	"eth-blockchain-service/internal/app"
	"eth-blockchain-service/internal/configs"
)

// @title           Ethereum Blockchain API Service
// @version         1.0
func main() {
	serverConfig := configs.GetConfig().Server
	server, err := app.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	port := fmt.Sprintf(":%s", serverConfig.Port)
	err = server.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}
