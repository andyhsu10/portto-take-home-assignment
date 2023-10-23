package app

import (
	"context"
	"log"
	"sync"
	"time"

	"eth-blockchain-service/internal/configs"
	"eth-blockchain-service/internal/services"
)

var recentBlockNum *uint64

func getRecentBlockNum() (num *uint64, err error) {
	if recentBlockNum == nil {
		ethClientSrv, err := services.NewEthClientService()
		if err != nil {
			return nil, err
		}

		ctx := context.Background()
		num, err := ethClientSrv.GetRecentBlockNum(ctx)
		if err != nil {
			return nil, err
		}

		recentBlockNum = num
	}
	return recentBlockNum, nil
}

func Run() error {
	blockSrv, err := services.NewBlockService()
	if err != nil {
		log.Fatalln("Failed to create new block service", err)
	}

	txnSrv, err := services.NewTxnService()
	if err != nil {
		log.Fatalln("Failed to create new txn service", err)
	}

	ethClientSrv, err := services.NewEthClientService()
	if err != nil {
		log.Fatalln("Failed to create new ETH client service", err)
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)

	// indexing past blocks
	go func() {
		defer wg.Done()
		num, err := getRecentBlockNum()
		if err != nil {
			log.Fatalln("Failed to get recent block number", err)
		}

		ctx := context.Background()
		blocks := makeRange(int(*num)-configs.GetConfig().MaxN, int(*num))
		blocksInDb, err := blockSrv.GetLatestNBlockNumbers(ctx, configs.GetConfig().MaxN)
		if err != nil {
			log.Fatalln("Failed to get latest block numbers from DB", err)
		}

		blocksInDbMap := make(map[int]bool)
		for _, b := range *blocksInDb {
			blocksInDbMap[b] = true
		}

		const rateLimit = 10
		ticker := time.NewTicker(time.Second / rateLimit)
		for _, b := range blocks {
			if _, ok := blocksInDbMap[b]; !ok {
				block, txns, err := ethClientSrv.GetBlock(ctx, uint64(b))
				if err != nil {
					log.Println("Failed to index the block:", b)
					continue
				}

				_, err = blockSrv.CreateBlock(ctx, *block)
				if err != nil {
					log.Println("Failed to create the block:", b)
					continue
				}

				_, err = txnSrv.BatchCreateTxns(ctx, *txns)
				if err != nil {
					log.Println("Failed to create the block's transactions:", b)
					continue
				}

				log.Println("Successfully index the past block:", b)

				<-ticker.C
			}
		}
	}()

	// indexing future blocks
	go func() {
		defer wg.Done()

		ticker := time.NewTicker(time.Second * 5)
		for {
			<-ticker.C

			ctx := context.Background()
			num, err := ethClientSrv.GetRecentBlockNum(ctx)
			if err != nil {
				continue
			}

			if recentBlockNum != nil && *num > *recentBlockNum {
				blocks := makeRange(int(*recentBlockNum)+1, int(*num))
				const rateLimit = 10
				ticker := time.NewTicker(time.Second / rateLimit)
				for _, b := range blocks {
					block, txns, err := ethClientSrv.GetBlock(ctx, uint64(b))
					if err != nil {
						log.Println("Failed to index the block:", b)
						continue
					}

					_, err = blockSrv.CreateBlock(ctx, *block)
					if err != nil {
						log.Println("Failed to create the block:", b)
						continue
					}

					_, err = txnSrv.BatchCreateTxns(ctx, *txns)
					if err != nil {
						log.Println("Failed to create the block's transactions:", b)
						continue
					}

					recentBlockNum = num
					log.Println("Successfully index the new block:", b)

					<-ticker.C
				}
			}
		}
	}()

	wg.Wait()

	return nil
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = max - i
	}
	return a
}
