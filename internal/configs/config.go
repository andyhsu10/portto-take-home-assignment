package configs

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	prod = "production"
)

var configInstance *Config

// Config object
type Config struct {
	Database *Database
	Env      string `env:"ENV"`
	Server   *Server
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

	return &Config{
		Database: &Database{
			URL: os.Getenv("DB_URL"),
		},
		Env: os.Getenv("ENV"),
		Server: &Server{
			Port: os.Getenv("SERVER_PORT"),
			Cors: strings.Split(os.Getenv("CORS_ORIGIN_WHITELIST"), ","),
		},
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
