# Portto Take-home assignment (Ethereum Blockchain Service)

## Prerequisites

- Go: v1.21.3
- Docker

## Getting Started

1. Install dependencies: `go mod download`
2. Copy `.env.example` to `.env` and fill in the values
   Example .env values:

   ```
   CORS_ORIGIN_WHITELIST="*"
   DB_DATABASE="ethereum"
   DB_PASSWORD="abcde12345"
   DB_URL="postgres://admin:abcde12345@localhost:5432/ethereum"
   DB_USER="admin"
   ENV="develop"
   MAX_N="10000"
   MAX_ROUTINES="5"
   RPC_LIST="https://data-seed-prebsc-2-s3.binance.org:8545/"
   SERVER_PORT="8080"
   ```

3. Start the PostgreSQL DB (in detached mode): `docker compose up -d`
4. Start the indexer service by: `go run ./cmd/indexer`
5. Start the API service by: `go run ./cmd/server`
6. To view the swagger UI, visit http://localhost:8080/swagger/index.html

## Requirements

1. API service

   - [GET] /blocks?limit=n (without transaction hash)
     Example response:
     ```
     {
       "blocks": [
         {
           "block_num": 1,
           "block_hash": "",
           "block_time": 123456789,
           "parent_hash": "",
         }
       ]
     }
     ```
   - [GET] /blocks/:id (with all transactions hashs)
     Example response:
     ```
     {
       "block_num": 1,
       "block_hash": "",
       "block_time": 123456789,
       "parent_hash": "",
       "transactions": [
         "0x12345678",
         "0x87654321"
       ]
     }
     ```
   - [GET] /transaction/:txHash
     Example response:
     ```
     {
       "tx_hash": "0x6666",
       "from": "0x4321",
       "to": "0x1234",
       "nonce": 1,
       "data": "0xeb12",
       "value": "12345678"
       "logs": [
         {
           "index": 0,
           "data": "0x12345678",
         }
       ]
     }
     ```

2. Ethereum block indexer service

## Miscellaneous

- Re-generate swagger docs: `swag init --g ./cmd/server/main.go`
