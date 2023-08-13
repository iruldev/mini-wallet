//go:build wireinject
// +build wireinject

package injector

import (
	_ "github.com/iruldev/mini-wallet/src/config"

	"github.com/iruldev/mini-wallet/engine/rest/transformer"
	"github.com/iruldev/mini-wallet/src/database"
	"github.com/iruldev/mini-wallet/src/repository"

	"github.com/go-playground/validator/v10"
	"github.com/iruldev/mini-wallet/engine/rest/controller"
	"github.com/iruldev/mini-wallet/src/service"
	"github.com/iruldev/mini-wallet/src/token"
	"github.com/spf13/viper"

	"github.com/google/wire"
)

func JwtSecretKey() string {
	return viper.GetString("JWT_SECRET_KEY")
}

func InitializeWalletService() service.WalletService {
	wire.Build(service.NewWalletService,
		token.NewJWTMaker,
		JwtSecretKey,
		repository.NewWalletRepository,
		database.GetDB,
	)
	return nil
}

func InitializeWalletControllerREST() (controller.WalletController, error) {
	wire.Build(
		controller.NewWalletController,
		validator.New,
		InitializeWalletService,
		transformer.NewWalletTransformer,
	)
	return nil, nil
}

func InitializeTransactionControllerREST() (controller.TransactionController, error) {
	wire.Build(
		controller.NewTransactionController,
		validator.New,
		service.NewTransactionService,
		transformer.NewTransactionTransformer,
		InitializeWalletService,
		repository.NewTransactionRepository,
		database.GetDB,
	)
	return nil, nil
}
