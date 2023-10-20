package services

// import (
// 	"gorm.io/gorm"
// )

type TxnService interface{}

type txnService struct {
	// db   *gorm.DB
}

func NewTxnService() (TxnService, error) {
	// db, err := databases.GetDB()
	// if err != nil {
	// 	return nil, err
	// }
	// return &blockService{db: db}, nil
	return &txnService{}, nil
}
