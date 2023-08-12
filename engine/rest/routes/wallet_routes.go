package routes

import (
	"net/http"

	injector "github.com/iruldev/mini-wallet/src"

	"github.com/iruldev/mini-wallet/src/middleware"
	"github.com/iruldev/mini-wallet/src/token"

	"github.com/gorilla/mux"
	"github.com/iruldev/mini-wallet/engine/rest/controller"
)

func WalletRoutes(r *mux.Router, c controller.WalletController) {
	r.HandleFunc("/init", c.InitWallet).Methods(http.MethodPost)
	walletR := r.PathPrefix("/wallet").Subrouter()
	walletR.Use(middleware.AuthMiddleware(token.NewJWTMaker(injector.JwtSecretKey())))
	walletR.HandleFunc("", c.EnableWallet).Methods(http.MethodPost)
	walletR.HandleFunc("", c.DisableWallet).Methods(http.MethodPatch)
	walletR.HandleFunc("", c.GetWallet).Methods(http.MethodGet)
}
