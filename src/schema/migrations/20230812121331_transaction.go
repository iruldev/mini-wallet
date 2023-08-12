package migrations

import (
	"context"
	"database/sql"
	"github.com/iruldev/mini-wallet/src/database"
	"github.com/iruldev/mini-wallet/src/model/entity"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upTransaction, downTransaction)
}

func upTransaction(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	db := database.GetDB()
	err := db.Set("gorm:table_options", TABLE_OPTIONS).AutoMigrate(&entity.Transaction{})
	if err != nil {
		return err
	}
	return nil
}

func downTransaction(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	db := database.GetDB()
	err := db.Migrator().DropTable(&entity.Transaction{})
	if err != nil {
		return err
	}
	return nil
}
