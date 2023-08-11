package controller

import "net/http"

type WalletController interface {
	InitWallet(w http.ResponseWriter, r *http.Request)
}
