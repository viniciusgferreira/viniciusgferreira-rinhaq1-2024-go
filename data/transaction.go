package data

import (
	"time"
)

type Transaction struct {
	Amount      int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao"`
	Date        time.Time `json:"realizada_em"`
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
