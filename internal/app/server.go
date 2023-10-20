package app

import (
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"

	"eth-blockchain-service/internal/routers"
)

type Server struct {
	apiRouter *gin.Engine
}

func NewServer() (*Server, error) {
	app, err := routers.GetRouter()
	if err != nil {
		return nil, err
	}

	return &Server{
		apiRouter: app,
	}, nil
}

func (server *Server) Run(port string) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Printf("server listen port:%s\n", port)
		log.Fatal(server.apiRouter.Run(port))
	}()
	wg.Wait()

	return nil
}
