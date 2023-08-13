package repository

import (
	"context"

	"github.com/iruldev/mini-wallet/src/model/entity"
)

type TransactionRepository interface {
	ForCustomerXID(customerXID string) TransactionRepository
	ForReferenceID(referenceID string) TransactionRepository
	ForProcessed() TransactionRepository

	Get(ctx context.Context) (*entity.Transaction, error)
	GetAll(ctx context.Context) ([]*entity.Transaction, error)
	Create(ctx context.Context, data *entity.Transaction) (*entity.Transaction, error)
	Update(ctx context.Context, data *entity.Transaction) error
	Process(ctx context.Context, data *entity.Transaction) error
	Complete(ctx context.Context, data *entity.Transaction) error
	Fail(ctx context.Context, data *entity.Transaction) error

	Clean() TransactionRepository
}
