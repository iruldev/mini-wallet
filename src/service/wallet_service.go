package service

import (
	"context"

	"github.com/iruldev/mini-wallet/src/model/entity"
)

type WalletService interface {
	InitWallet(ctx context.Context, customerXID string) (string, error)
	ActivateWallet(ctx context.Context, customerXID string) (*entity.Wallet, error)
	DeactivateWallet(ctx context.Context, customerXID string) (*entity.Wallet, error)
	GetWallet(ctx context.Context, customerXID string) (*entity.Wallet, error)
}
