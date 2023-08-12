package repository

import (
	"context"
	"errors"

	"github.com/iruldev/mini-wallet/src/model/entity"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
	db         *gorm.DB
	whereQuery *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &TransactionRepositoryImpl{db: db}
}

func (r *TransactionRepositoryImpl) buildQuery() *gorm.DB {
	if r.whereQuery != nil {
		return r.whereQuery
	}
	return r.db
}

func (r *TransactionRepositoryImpl) ForReferenceID(referenceID string) TransactionRepository {
	r.whereQuery = r.buildQuery().Where("reference_id = ?", referenceID)
	return r
}

func (r *TransactionRepositoryImpl) Get(ctx context.Context) (*entity.Transaction, error) {
	defer r.Clean() // clean where query on executed
	if r.whereQuery == nil {
		return nil, errors.New("filter query is required")
	}
	var data entity.Transaction
	err := r.whereQuery.First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &data, nil
}

func (r *TransactionRepositoryImpl) GetAll(ctx context.Context) ([]*entity.Transaction, error) {
	defer r.Clean() // clean where query on executed
	db := r.buildQuery().WithContext(ctx)

	var funds []*entity.Transaction
	err := db.Find(&funds).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return funds, nil
		}
		return nil, err
	}

	return funds, nil
}

func (r *TransactionRepositoryImpl) Create(ctx context.Context, data *entity.Transaction) (*entity.Transaction, error) {
	defer r.Clean() // clean where query on executed
	if data == nil {
		return nil, errors.New("data is required")
	}
	err := r.db.WithContext(ctx).Create(data).Error
	return data, err
}

func (r *TransactionRepositoryImpl) Update(ctx context.Context, data *entity.Transaction) error {
	defer r.Clean() // clean where query on executed
	if data == nil {
		return errors.New("data is required")
	}

	db := r.db.WithContext(ctx)

	if db.First(&entity.Transaction{}, data.ID).RowsAffected == 0 {
		return errors.New("transaction not found")
	}

	return db.Save(data).Error
}

func (r *TransactionRepositoryImpl) Clean() TransactionRepository {
	r.whereQuery = nil
	return r
}
