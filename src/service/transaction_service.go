package service

import (
	"context"

	"github.com/iruldev/mini-wallet/src/model/entity"
)

type TransactionService interface {
	GetTransactions(ctx context.Context, customerXID string) ([]*entity.Transaction, error)
	Transaction(ctx context.Context, customerXID string, params TransactionReq) (*entity.Transaction, error)
	CompleteTransaction(ctx context.Context, data entity.Transaction) (*entity.Transaction, error)
}
