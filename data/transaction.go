package data

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

type Transaction struct {
	Amount      int        `json:"valor"`
	Type        string     `json:"tipo"`
	Description string     `json:"descricao"`
	Date        *time.Time `json:"realizada_em"`
}

func NewTransaction(r io.Reader) (*Transaction, error) {
	t := &Transaction{}
	err := t.FromJSON(r)
	if err != nil {
		return nil, errors.New("")
	}
	now := time.Now()
	t.Date = &now
	return t, nil
}

func (t *Transaction) Validate() bool {
	if t.Amount == 0 || len(t.Type) <= 0 || len(t.Description) <= 0 {
		return false
	}
	if "c" != t.Type && t.Type != "d" {
		return false
	}

	if len(t.Description) > 10 {
		return false
	}

	return true
}
func (t *Transaction) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func (t *Transaction) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(t)
}
