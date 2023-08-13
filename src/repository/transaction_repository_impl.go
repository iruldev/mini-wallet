package repository

import (
	"context"
	"errors"
	"fmt"

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

func (r *TransactionRepositoryImpl) ForCustomerXID(customerXID string) TransactionRepository {
	r.whereQuery = r.buildQuery().Where("customer_xid = ?", customerXID)
	return r
}

func (r *TransactionRepositoryImpl) ForReferenceID(referenceID string) TransactionRepository {
	r.whereQuery = r.buildQuery().Where("reference_id = ?", referenceID)
	return r
}

func (r *TransactionRepositoryImpl) ForProcessed() TransactionRepository {
	r.whereQuery = r.buildQuery().Where("is_processed = 1")
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

	var trns []*entity.Transaction
	err := r.buildQuery().Order("created_at DESC").Find(&trns).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return trns, nil
		}
		return nil, err
	}

	return trns, nil
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

func (r *TransactionRepositoryImpl) Complete(ctx context.Context, data *entity.Transaction) error {
	defer r.Clean() // clean where query on executed
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Recovered in e", e)
		}
	}()

	db := r.db.WithContext(ctx)
	// CHECK WITHOUT LOCKING
	if db.First(&data, "id = ?", data.ID).RowsAffected == 0 {
		return errors.New("transaction not found")
	}

	if data.IsCompleted > 0 {
		return errors.New("transaction already complete")
	}

	if data.IsFailed > 0 {
		return errors.New("transaction already failed")
	}

	if data.IsProcessed == 0 {
		e := r.Process(ctx, data)
		if e != nil {
			return e
		}
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// DOUBLE CHECK WITH LOCKING
		if tx.First(&data, "id = ?", data.ID).RowsAffected == 0 {
			return errors.New("transaction not found")
		}

		if data.IsProcessed == 0 {
			return errors.New("transaction not processed")
		}

		if data.IsCompleted > 0 {
			return errors.New("transaction already complete")
		}

		if data.IsFailed > 0 {
			return errors.New("transaction already failed")
		}

		var wallet *entity.Wallet
		if tx.First(&wallet, "customer_xid = ?", data.CustomerXID).RowsAffected == 0 {
			return errors.New("wallet not found")
		}

		if data.Type == entity.WITHDRAWAL {
			if data.Amount.GreaterThan(wallet.Balance) {
				return errors.New("balance not enough")
			}
			wallet.Balance = wallet.Balance.Sub(data.Amount)
		} else if data.Type == entity.DEPOSIT {
			wallet.Balance = wallet.Balance.Add(data.Amount)
		} else {
			return errors.New("unsupported type")
		}

		if err := tx.Save(wallet).Error; err != nil {
			return err
		}

		data.IsCompleted = 1
		if err := tx.Save(data).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		_ = r.Fail(ctx, data)
	}

	return err
}

func (r *TransactionRepositoryImpl) Process(ctx context.Context, data *entity.Transaction) error {
	defer r.Clean() // clean where query on executed
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Recovered in e", e)
		}
	}()

	db := r.db.WithContext(ctx)
	err := db.Transaction(func(tx *gorm.DB) error {
		// DOUBLE CHECK WITH LOCKING
		if tx.First(&data, "id = ?", data.ID).RowsAffected == 0 {
			return errors.New("transaction not found")
		}

		if data.IsProcessed > 0 {
			return errors.New("transaction already processed")
		}

		if data.IsCompleted > 0 {
			return errors.New("transaction already complete")
		}

		if data.IsFailed > 0 {
			return errors.New("transaction already failed")
		}

		data.IsProcessed = 1
		if err := tx.Save(data).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		_ = r.Fail(ctx, data)
	}

	return err
}

func (r *TransactionRepositoryImpl) Fail(ctx context.Context, data *entity.Transaction) error {
	defer r.Clean() // clean where query on executed
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("Recovered in e", e)
		}
	}()

	db := r.db.WithContext(ctx)
	err := db.Transaction(func(tx *gorm.DB) error {
		// DOUBLE CHECK WITH LOCKING
		if tx.First(&data, "id = ?", data.ID).RowsAffected == 0 {
			return errors.New("transaction not found")
		}

		if data.IsCompleted > 0 {
			return errors.New("transaction already complete")
		}

		if data.IsFailed > 0 {
			return errors.New("transaction already fail")
		}

		data.IsFailed = 1
		if err := tx.Save(data).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func (r *TransactionRepositoryImpl) Clean() TransactionRepository {
	r.whereQuery = nil
	return r
}
