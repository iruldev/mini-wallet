package service

import "context"

type WalletService interface {
	InitWallet(ctx context.Context, params InitWalletReq) (string, error)
}
