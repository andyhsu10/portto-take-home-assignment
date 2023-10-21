package models

import "time"

type BaseModel struct {
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Block struct {
	BaseModel
	Number       int           `json:"block_num" gorm:"primaryKey;uniqueIndex:,sort:desc"`
	Hash         string        `json:"block_hash"`
	Time         int           `json:"block_time"`
	ParentHash   string        `json:"parent_hash"`
	Transactions []Transaction `json:"transactions" gorm:"foreignKey:BlockNumber;references:Number"`
}

type Transaction struct {
	BaseModel
	Hash        string `json:"tx_hash" gorm:"primaryKey;type:char(66);uniqueIndex"`
	From        string `json:"from"`
	To          string `json:"to"`
	Nonce       int    `json:"nonce"`
	Data        string `json:"data"`
	Value       string `json:"value"`
	Logs        string `json:"logs"`
	BlockNumber int    `json:"-"`
}
