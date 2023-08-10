package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iruldev/mini-wallet/src/constant"
	"github.com/iruldev/mini-wallet/src/helper"
)

func AppRoutes(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		res := helper.PlugResponse(w)
		_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, "Service is running well"))
	}).Methods(http.MethodGet)
}
