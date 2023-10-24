package services

import (
	"context"
	"errors"
	"eth-blockchain-service/internal/configs"
	"log"
	"time"
)

type IndexerService interface {
	IndexOldBlocks(ctx context.Context, begin int, block chan int) error
	IndexNewBlocks(ctx context.Context, block chan int) error
	IndexOneBlock(ctx context.Context, blockNum uint64) error
}

type indexerService struct {
	recentBlockNum uint64
	blockSrv       BlockService
	txnSrv         TxnService
	ethClientSrv   EthClientService
}

func NewIndexerService() (IndexerService, error) {
	blockSrv, err := NewBlockService()
	if err != nil {
		log.Fatalln("Failed to create new block service", err)
		return nil, err
	}

	txnSrv, err := NewTxnService()
	if err != nil {
		log.Fatalln("Failed to create new txn service", err)
		return nil, err
	}

	ethClientSrv, err := NewEthClientService()
	if err != nil {
		log.Fatalln("Failed to create new ETH client service", err)
		return nil, err
	}

	ctx := context.Background()
	number, err := ethClientSrv.GetRecentBlockNum(ctx)
	if err != nil {
		log.Fatalln("Failed to get recent block number", err)
		return nil, err
	}

	return &indexerService{
		recentBlockNum: number,
		blockSrv:       blockSrv,
		txnSrv:         txnSrv,
		ethClientSrv:   ethClientSrv,
	}, nil
}

func (srv *indexerService) IndexOldBlocks(ctx context.Context, begin int, block chan int) error {
	if begin <= 0 || begin > int(srv.recentBlockNum) {
		close(block)
		return errors.New("begin number must be greater than zero or it is greater than the recent block number")
	}

	count := int(srv.recentBlockNum) - begin + 1
	chunkSize := 1000
	iter := count/chunkSize + 1
	for i := 0; i < iter; i++ {
		start := begin + i*chunkSize
		var blocks []int
		if i < iter-1 {
			blocks = srv.makeRange(start, start+chunkSize-1)
		} else {
			blocks = srv.makeRange(start, int(srv.recentBlockNum))
		}

		blocksInDb, err := srv.blockSrv.GetBlockNumbers(ctx, &blocks)
		if err != nil {
			log.Fatalln("Failed to get latest block numbers from DB", err)
			return err
		}

		blocksInDbMap := make(map[int]bool)
		for _, b := range *blocksInDb {
			blocksInDbMap[b] = true
		}

		ticker := time.NewTicker(time.Second / time.Duration(configs.GetConfig().IndexRateLimit))
		for _, b := range blocks {
			if _, ok := blocksInDbMap[b]; !ok {
				err := srv.IndexOneBlock(ctx, uint64(b))
				if err == nil {
					block <- b
				}

				<-ticker.C
			}
		}
	}

	close(block)
	return nil
}

func (srv *indexerService) IndexNewBlocks(ctx context.Context, block chan int) error {
	ticker := time.NewTicker(time.Second * 5)
	for {
		<-ticker.C

		num, err := srv.ethClientSrv.GetRecentBlockNum(ctx)
		if err != nil {
			continue
		}

		if num > srv.recentBlockNum {
			blocks := srv.makeRange(int(srv.recentBlockNum)+1, int(num))
			ticker := time.NewTicker(time.Second / time.Duration(configs.GetConfig().IndexRateLimit))

			for _, b := range blocks {
				err := srv.IndexOneBlock(ctx, uint64(b))
				if err != nil {
					break
				}

				srv.recentBlockNum = uint64(b)
				block <- b
				<-ticker.C
			}
		}
	}
}

func (srv *indexerService) IndexOneBlock(ctx context.Context, blockNum uint64) error {
	block, txns, err := srv.ethClientSrv.GetBlock(ctx, blockNum)
	if err != nil {
		log.Println("Failed to index the block:", blockNum)
		return err
	}

	_, err = srv.blockSrv.CreateBlock(ctx, *block)
	if err != nil {
		log.Println("Failed to create the block:", blockNum)
		return err
	}

	_, err = srv.txnSrv.BatchCreateTxns(ctx, *txns)
	if err != nil {
		log.Println("Failed to create the block's transactions:", blockNum)
		return err
	}

	return nil
}

func (srv *indexerService) makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}
