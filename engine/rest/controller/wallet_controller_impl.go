package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/iruldev/mini-wallet/src/constant"
	"github.com/iruldev/mini-wallet/src/helper"
	"github.com/iruldev/mini-wallet/src/service"
	"net/http"
)

type WalletControllerImpl struct {
	Validator *validator.Validate
}

func NewWalletController(validator *validator.Validate) WalletController {
	return &WalletControllerImpl{Validator: validator}
}

func (c WalletControllerImpl) InitWallet(w http.ResponseWriter, r *http.Request) {
	req := helper.PlugRequest(r, w)
	res := helper.PlugResponse(w)

	pReq, _ := helper.ParseTo[service.InitWalletReq](req)
	err := c.Validator.Struct(pReq)
	if err != nil {
		errF := helper.GetErrMsgField(err)
		_ = res.ReplyCustom(http.StatusBadRequest, helper.NewResponse(constant.FAILED, errF))
		return
	}

	_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, struct{}{}))
}
