package data

import (
	"time"
)

type Client struct {
	ID           uint8          `json:"id,omitempty"`
	Balance      *Balance       `json:"saldo"`
	Transactions []*Transaction `json:"ultimas_transacoes,omitempty"`
}

type Balance struct {
	Total int       `json:"total"`
	Limit int       `json:"limite"`
	Date  time.Time `json:"data_extrato,omitempty"`
}

type ResponseTransaction struct {
	Limit int `json:"limite"`
	Total int `json:"saldo"`
}
