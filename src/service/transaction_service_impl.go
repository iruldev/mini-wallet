package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/iruldev/mini-wallet/src/model/entity"
	"github.com/iruldev/mini-wallet/src/repository"
	"github.com/shopspring/decimal"
)

type TransactionServiceImpl struct {
	Repository repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &TransactionServiceImpl{Repository: repo}
}

func (s *TransactionServiceImpl) GetTransactions(ctx context.Context, customerXID string) ([]*entity.Transaction, error) {
	return s.Repository.ForCustomerXID(customerXID).ForProcessed().GetAll(ctx)
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

	if amount.IsZero() || amount.IsNegative() {
		return nil, errors.New(fmt.Sprintf("amount must be greater than 0"))
	}

	trns, err := s.Repository.Create(ctx, &entity.Transaction{
		CustomerXID: customerXID,
		Type:        params.Type,
		ReferenceID: params.ReferenceID,
		At:          &nwTime,
		Amount:      amount,
	})

	if err != nil {
		return nil, err
	}

	trx, err := s.CompleteTransaction(ctx, *trns)
	if err != nil {
		return nil, err
	}

	return trx, nil
}

func (s *TransactionServiceImpl) CompleteTransaction(ctx context.Context, data entity.Transaction) (*entity.Transaction, error) {
	if data.IsCompleted > 0 {
		return nil, errors.New("transaction already complete")
	}

	if data.IsFailed > 0 {
		return nil, errors.New("transaction already fail")
	}

	e := s.Repository.Complete(ctx, &data)
	if e != nil {
		return nil, e
	}

	return &data, nil
}
