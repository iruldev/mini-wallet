package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iruldev/mini-wallet/engine/rest/controller"
	injector "github.com/iruldev/mini-wallet/src"
	"github.com/iruldev/mini-wallet/src/middleware"
	"github.com/iruldev/mini-wallet/src/token"
)

func TransactionRoutes(r *mux.Router, c controller.TransactionController) {
	transactionR := r.PathPrefix("/wallet").Subrouter()
	transactionR.Use(middleware.AuthMiddleware(token.NewJWTMaker(injector.JwtSecretKey())))
	transactionR.HandleFunc("/transactions", c.GetTransactions).Methods(http.MethodGet)
	transactionR.HandleFunc("/deposits", c.Deposit).Methods(http.MethodPost)
	transactionR.HandleFunc("/withdrawals", c.Withdrawal).Methods(http.MethodPost)
}
