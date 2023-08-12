package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/iruldev/mini-wallet/src/model/entity"
	"github.com/iruldev/mini-wallet/src/repository"

	"github.com/iruldev/mini-wallet/src/token"
)

type WalletServiceImpl struct {
	Token      token.Maker
	Repository repository.WalletRepository
}

func NewWalletService(token token.Maker, repo repository.WalletRepository) WalletService {
	return &WalletServiceImpl{
		Token:      token,
		Repository: repo,
	}
}

func (s *WalletServiceImpl) InitWallet(ctx context.Context, customerXID string) (string, error) {
	wallet, err := s.Repository.ForCustomerXID(customerXID).Get(ctx)
	if err != nil || wallet == nil {
		err = s.Repository.Create(ctx, &entity.Wallet{
			CustomerXID: customerXID,
		})
		if err != nil {
			return "", err
		}
	}
	token, _, err := s.Token.CreateToken(customerXID, time.Minute*10)
	return token, err
}

func (s *WalletServiceImpl) ActivateWallet(ctx context.Context, customerXID string) (*entity.Wallet, error) {
	wallet, err := s.Repository.ForCustomerXID(customerXID).Get(ctx)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, errors.New(fmt.Sprintf("Wallet %s not found", customerXID))
	}

	if wallet.IsActive == 1 {
		return nil, errors.New("Already enabled")
	}

	err = s.Repository.ForCustomerXID(wallet.CustomerXID).Activate(ctx)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletServiceImpl) DeactivateWallet(ctx context.Context, customerXID string) (*entity.Wallet, error) {
	wallet, err := s.Repository.ForCustomerXID(customerXID).Get(ctx)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, errors.New(fmt.Sprintf("Wallet %s not found", customerXID))
	}

	err = s.Repository.ForCustomerXID(wallet.CustomerXID).Deactivate(ctx)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletServiceImpl) GetWallet(ctx context.Context, customerXID string) (*entity.Wallet, error) {
	wallet, err := s.Repository.ForCustomerXID(customerXID).Get(ctx)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		return nil, errors.New(fmt.Sprintf("Wallet %s not found", customerXID))
	}

	if wallet.IsActive == 0 {
		return nil, errors.New("Wallet disabled")
	}

	return wallet, nil
}
