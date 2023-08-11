package migrations

import (
	"context"
	"database/sql"

	"github.com/iruldev/mini-wallet/src/database"
	"github.com/iruldev/mini-wallet/src/model/entity"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upWallet, downWallet)
}

func upWallet(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	db := database.GetDB()
	err := db.Set("gorm:table_options", TABLE_OPTIONS).AutoMigrate(&entity.Wallet{})
	if err != nil {
		return err
	}
	return nil
}

func downWallet(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	db := database.GetDB()
	err := db.Migrator().DropTable(&entity.Wallet{})
	if err != nil {
		return err
	}
	return nil
}
