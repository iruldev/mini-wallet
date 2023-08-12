package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/iruldev/mini-wallet/src/model/entity"
	"github.com/iruldev/mini-wallet/src/repository"
	"github.com/shopspring/decimal"
	"time"
)

type TransactionServiceImpl struct {
	Repository repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &TransactionServiceImpl{Repository: repo}
}

func (s *TransactionServiceImpl) GetTransactions(ctx context.Context, customerXID string) ([]*entity.Transaction, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TransactionServiceImpl) Transaction(ctx context.Context, customerXID string, params TransactionReq) (*entity.Transaction, error) {
	t, err := s.Repository.ForReferenceID(params.ReferenceID).Get(ctx)
	if err != nil {
		return nil, err
	}

	if t != nil {
		return nil, errors.New(fmt.Sprintf("transaction with reference %s already exist", params.ReferenceID))
	}

	nwTime := time.Now()
	amount, err := decimal.NewFromString(params.Amount)
	if err != nil {
		return nil, err
	}

	trns, err := s.Repository.Create(ctx, &entity.Transaction{
		CustomerXID: customerXID,
		Type:        params.Type,
		ReferenceID: params.ReferenceID,
		At:          &nwTime,
		Amount:      amount,
	})

	return trns, nil
}
