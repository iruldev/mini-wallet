//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/go-playground/validator/v10"
	"github.com/iruldev/mini-wallet/engine/rest/controller"

	"github.com/google/wire"
)

func InitializeWalletControllerREST() (controller.WalletController, error) {
	wire.Build(controller.NewWalletController, validator.New)
	return nil, nil
}
