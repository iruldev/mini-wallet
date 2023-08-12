package controller

import "net/http"

type WalletController interface {
	InitWallet(w http.ResponseWriter, r *http.Request)
	EnableWallet(w http.ResponseWriter, r *http.Request)
	DisableWallet(w http.ResponseWriter, r *http.Request)
	GetWallet(w http.ResponseWriter, r *http.Request)
}
