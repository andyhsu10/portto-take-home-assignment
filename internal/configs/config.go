package configs

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	prod = "production"
)

var configInstance *Config

// Config object
type Config struct {
	Database              *Database
	Env                   string
	IndexRateLimit        int
	MaxN                  int
	RpcList               []string
	Server                *Server
	StartIndexBlockNumber int
}

type Database struct {
	URL string
}

type Server struct {
	Port string
	Cors []string
}

// IsProd Checks if env is production
func (c Config) IsProd() bool {
	return c.Env == prod
}

func newConfig() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	} else if os.IsNotExist(err) {
		log.Println("no .env file")
	}

	indexRateLimit, err := strconv.Atoi(os.Getenv("INDEX_RATE_LIMIT"))
	if err != nil || indexRateLimit <= 0 {
		indexRateLimit = 15
	}

	maxN, err := strconv.Atoi(os.Getenv("MAX_N"))
	if err != nil || maxN <= 0 {
		maxN = 10000
	}

	startIndexBlockNumber, err := strconv.Atoi(os.Getenv("START_INDEX_BLOCK_NUMBER"))
	if err != nil || startIndexBlockNumber < 0 {
		startIndexBlockNumber = 0
	}

	return &Config{
		Database: &Database{
			URL: os.Getenv("DB_URL"),
		},
		Env:            os.Getenv("ENV"),
		IndexRateLimit: indexRateLimit,
		MaxN:           maxN,
		RpcList:        strings.Split(os.Getenv("RPC_LIST"), ","),
		Server: &Server{
			Port: os.Getenv("SERVER_PORT"),
			Cors: strings.Split(os.Getenv("CORS_ORIGIN_WHITELIST"), ","),
		},
		StartIndexBlockNumber: startIndexBlockNumber,
	}, nil
}

// GetConfig gets all config for the application
func GetConfig() *Config {
	if configInstance == nil {
		instance, err := newConfig()
		if err != nil {
			log.Fatal(err)
		}
		configInstance = instance
	}

	return configInstance
}
