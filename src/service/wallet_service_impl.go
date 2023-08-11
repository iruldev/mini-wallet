package service

import (
	"context"
	"time"

	"github.com/iruldev/mini-wallet/src/token"
)

type WalletServiceImpl struct {
	Token token.Maker
}

func NewWalletService(token token.Maker) WalletService {
	return &WalletServiceImpl{Token: token}
}

func (s *WalletServiceImpl) InitWallet(ctx context.Context, params InitWalletReq) (string, error) {
	token, _, err := s.Token.CreateToken(params.CustomerXID, time.Minute*10)
	return token, err
}
