package databases

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"eth-blockchain-service/internal/configs"
)

// set up a singleton instance of the database
var dbInstance *gorm.DB

func GetDB() (instance *gorm.DB, err error) {
	if dbInstance == nil {
		instance, err = newDB()
		if err != nil {
			return nil, err
		}
		dbInstance = instance
	}
	return dbInstance, nil
}

func newDB() (*gorm.DB, error) {
	dbConfig := configs.GetConfig().Database
	db, err := gorm.Open(
		postgres.Open(dbConfig.URL))
	if err != nil {
		return nil, err
	}

	// TODO: Auto-migrate models

	return db, nil
}
