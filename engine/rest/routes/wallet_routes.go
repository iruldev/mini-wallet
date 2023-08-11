package routes

import (
	"github.com/gorilla/mux"
	"github.com/iruldev/mini-wallet/engine/rest/controller"
	"net/http"
)

func WalletRoutes(r *mux.Router, c controller.WalletController) {
	r.HandleFunc("/init", c.InitWallet).Methods(http.MethodPost)
}
