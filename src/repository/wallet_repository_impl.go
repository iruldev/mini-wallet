package repository

import (
	"context"
	"errors"
	"time"

	"github.com/iruldev/mini-wallet/src/model/entity"
	"gorm.io/gorm"
)

type WalletRepositoryImpl struct {
	db         *gorm.DB
	whereQuery *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &WalletRepositoryImpl{
		db: db,
	}
}

func (r *WalletRepositoryImpl) buildQuery() *gorm.DB {
	if r.whereQuery != nil {
		return r.whereQuery
	}
	return r.db
}

func (r *WalletRepositoryImpl) ForCustomerXID(customerXID string) WalletRepository {
	r.whereQuery = r.buildQuery().Where("customer_xid = ?", customerXID)
	return r
}

func (r *WalletRepositoryImpl) Get(ctx context.Context) (*entity.Wallet, error) {
	defer r.Clean() // clean where query on executed
	if r.whereQuery == nil {
		return nil, errors.New("filter query is required")
	}
	var data entity.Wallet
	err := r.whereQuery.First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (r *WalletRepositoryImpl) Create(ctx context.Context, data *entity.Wallet) error {
	defer r.Clean() // clean where query on executed
	if data == nil {
		return errors.New("data is required")
	}
	return r.db.WithContext(ctx).Create(data).Error
}

func (r *WalletRepositoryImpl) Update(ctx context.Context, data *entity.Wallet) error {
	defer r.Clean() // clean where query on executed
	if data == nil {
		return errors.New("data is required")
	}

	db := r.db.WithContext(ctx)

	if db.First(&entity.Wallet{}, data.ID).RowsAffected == 0 {
		return errors.New("wallet not found")
	}

	return db.Save(data).Error
}

func (r *WalletRepositoryImpl) Activate(ctx context.Context) error {
	defer r.Clean() // clean where query on executed
	if r.whereQuery == nil {
		return errors.New("filter query is required")
	}

	db := r.buildQuery().WithContext(ctx)

	var data *entity.Wallet
	if db.Find(&data).RowsAffected == 0 {
		return errors.New("no wallet found")
	}

	data.IsActive = 1
	nwTime := time.Now()
	data.ActivatedAt = &nwTime

	return db.Save(data).Error
}

func (r *WalletRepositoryImpl) Deactivate(ctx context.Context) error {
	defer r.Clean() // clean where query on executed
	if r.whereQuery == nil {
		return errors.New("filter query is required")
	}

	db := r.buildQuery().WithContext(ctx)

	var data *entity.Wallet
	if db.Find(&data).RowsAffected == 0 {
		return errors.New("no wallet found")
	}

	data.IsActive = 0
	data.ActivatedAt = nil

	return db.Save(data).Error
}

func (r *WalletRepositoryImpl) Clean() WalletRepository {
	r.whereQuery = nil
	return r
}
