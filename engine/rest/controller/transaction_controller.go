package controller

import "net/http"

type TransactionController interface {
	GetTransactions(w http.ResponseWriter, r *http.Request)
	Deposit(w http.ResponseWriter, r *http.Request)
	Withdrawal(w http.ResponseWriter, r *http.Request)
}
