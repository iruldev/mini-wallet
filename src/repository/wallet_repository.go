package repository

import (
	"context"

	"github.com/iruldev/mini-wallet/src/model/entity"
)

type WalletRepository interface {
	ForCustomerXID(customerXID string) WalletRepository

	Get(ctx context.Context) (*entity.Wallet, error)
	Create(ctx context.Context, data *entity.Wallet) error
	Update(ctx context.Context, data *entity.Wallet) error
	Activate(ctx context.Context) error
	Deactivate(ctx context.Context) error

	Clean() WalletRepository
}
