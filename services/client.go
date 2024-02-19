package services

import (
	"errors"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/data"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/repositories"
)

type Service struct {
	repo *repositories.Repository
}

func New(repo *repositories.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTransaction(t *data.Transaction, id uint8) (*data.Balance, error) {
	var balance *data.Balance
	var err error
	if t.Type == "d" {
		balance, err = s.repo.Debit(t, id)
		if err != nil {
			return nil, errors.New("not enough funds or limit")
		}
	} else {
		balance, err = s.repo.Credit(t, id)
		if err != nil {
			return nil, err
		}
	}
	return balance, nil
}
func (s *Service) SufficientFunds(id uint8, amount int) bool {
	balance := s.repo.FindBalance(id)
	return amount <= balance.Total+balance.Limit
}

func (s *Service) CreateStatement(id uint8) *data.Client {
	client := s.repo.CreateStatement(id)
	return client
}
