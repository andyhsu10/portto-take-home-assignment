package main

import (
	"fmt"
	"log"

	"eth-blockchain-service/internal/app"
	"eth-blockchain-service/internal/configs"
)

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
