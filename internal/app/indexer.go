package app

import (
	"context"
	"log"
	"sync"

	"eth-blockchain-service/internal/configs"
	"eth-blockchain-service/internal/services"
)

func Run() error {
	indexSrv, err := services.NewIndexerService()
	if err != nil {
		log.Fatalln("Failed to get new indexer service", err)
		return err
	}

	ctx := context.Background()
	oldBlockNum := make(chan int)
	newBlockNum := make(chan int)
	wg := &sync.WaitGroup{}

	wg.Add(4)
	// indexing old blocks
	go func() {
		defer wg.Done()
		begin := configs.GetConfig().StartIndexBlockNumber
		if begin <= 0 {
			ethClientSrv, err := services.NewEthClientService()
			if err == nil {
				number, err := ethClientSrv.GetRecentBlockNum(ctx)
				if err == nil {
					begin = int(number) - 10000
				}
			}
		}
		indexSrv.IndexOldBlocks(ctx, begin, oldBlockNum)
	}()

	// indexing future new blocks
	go func() {
		defer wg.Done()
		indexSrv.IndexNewBlocks(ctx, newBlockNum)
	}()

	go func() {
		defer wg.Done()
		for num := range oldBlockNum {
			log.Println("Successfully index the old block:", num)
		}
	}()

	go func() {
		defer wg.Done()
		for num := range newBlockNum {
			log.Println("Successfully index the new block:", num)
		}
	}()

	wg.Wait()

	return nil
}
