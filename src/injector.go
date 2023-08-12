//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/iruldev/mini-wallet/engine/rest/transformer"
	_ "github.com/iruldev/mini-wallet/src/config"
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

func InitializeWalletControllerREST() (controller.WalletController, error) {
	wire.Build(
		controller.NewWalletController,
		validator.New,
		service.NewWalletService,
		token.NewJWTMaker,
		JwtSecretKey,
		transformer.NewWalletTransformer,
		repository.NewWalletRepository,
		database.GetDB,
	)
	return nil, nil
}
