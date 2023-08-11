package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iruldev/mini-wallet/engine/rest/controller"
)

func WalletRoutes(r *mux.Router, c controller.WalletController) {
	r.HandleFunc("/init", c.InitWallet).Methods(http.MethodPost)
}
