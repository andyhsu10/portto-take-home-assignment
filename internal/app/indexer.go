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
		indexSrv.IndexOldBlocks(ctx, configs.GetConfig().StartIndexingBlockNumber, oldBlockNum)
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
