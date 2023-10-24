# Portto Take-home assignment (Ethereum Blockchain Service)

## Prerequisites

- Go: v1.21.3
- Docker

## Getting Started

1. Clone the repo: `git clone git@github.com:andyhsu10/portto-take-home-assignment.git`
2. Install dependencies: `go mod download`
3. Copy `.env.example` to `.env` and fill in the values
   Example .env values:

   ```
   CORS_ORIGIN_WHITELIST="*"
   DB_DATABASE="ethereum"
   DB_PASSWORD="abcde12345"
   DB_URL="postgres://admin:abcde12345@localhost:5432/ethereum"
   DB_USER="admin"
   ENV="develop"
   INDEX_RATE_LIMIT="15"
   MAX_N="10000"
   RPC_LIST="https://data-seed-prebsc-2-s3.binance.org:8545/"
   SERVER_PORT="8080"
   START_INDEX_BLOCK_NUMBER="34500000"
   ```

4. Start the PostgreSQL DB (in detached mode): `docker compose up -d`
5. Start the indexer service by: `go run ./cmd/indexer`
6. Open another terminal and start the API service by: `go run ./cmd/server`
7. DB schema is located at: `/internal/models/model.go`
8. To view the swagger UI, visit http://localhost:8080/swagger/index.html

## Requirements

1. API service

   - [GET] /blocks?limit=n (without transaction hash)
     Description: Returns the latest n blocks from the DB
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
     Description: Returns the block from the DB if it exist
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
     Description: Returns the transaction from the DB if it exist. Logs will be fetched by the RPC once the corresponding field in DB is empty.
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
