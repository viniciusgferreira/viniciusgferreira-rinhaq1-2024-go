package services

import (
	"context"
	"errors"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/data"
	"github.com/viniciusgferreira/viniciusgferreira-rinhaq1-2024-go/repositories"
	"time"
)

type Service struct {
	repo *repositories.Repository
}

var (
	ErrInsufficientFunds = errors.New("not enough funds or limit")
)

func New(repo *repositories.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTransaction(ctx context.Context, t *data.Transaction, id int) (*data.ResponseTransaction, error) {
	t.Date = time.Now()
	var response *data.ResponseTransaction
	var err error
	if t.Type == "d" {
		response, err = s.repo.Debit(ctx, t, id)
		if err != nil {
			return nil, ErrInsufficientFunds
		}
	} else if t.Type == "c" {
		response, err = s.repo.Credit(ctx, t, id)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("wrong type (c,d)")
	}
	return response, nil
}

func (s *Service) CreateStatement(ctx context.Context, id int) *data.Client {
	client := s.repo.CreateStatement(ctx, id)
	return client
}
