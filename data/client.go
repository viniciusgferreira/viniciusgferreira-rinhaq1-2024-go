package data

import (
	"encoding/json"
	"io"
	"time"
)

type Client struct {
	ID           uint8          `json:"id,omitempty"`
	Balance      *Balance       `json:"saldo"`
	Transactions []*Transaction `json:"ultimas_transacoes,omitempty"`
}

type Balance struct {
	Total int        `json:"total"`
	Limit int        `json:"limite"`
	Date  *time.Time `json:"data_extrato,omitempty"`
}

func (t *Balance) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (t *Balance) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}

func (c *Client) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(c)
}

func (c *Client) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(c)
}
