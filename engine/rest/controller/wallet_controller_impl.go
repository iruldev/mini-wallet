package controller

import (
	"net/http"

	"github.com/iruldev/mini-wallet/src/token"

	"github.com/iruldev/mini-wallet/engine/rest/transformer"

	"github.com/go-playground/validator/v10"
	"github.com/iruldev/mini-wallet/src/constant"
	"github.com/iruldev/mini-wallet/src/helper"
	"github.com/iruldev/mini-wallet/src/service"
)

type WalletControllerImpl struct {
	Validator   *validator.Validate
	Service     service.WalletService
	Transformer transformer.WalletTransformer
}

func NewWalletController(
	validator *validator.Validate,
	walletService service.WalletService,
	transformer transformer.WalletTransformer,
) WalletController {
	return &WalletControllerImpl{
		Validator:   validator,
		Service:     walletService,
		Transformer: transformer,
	}
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

	token, _ := c.Service.InitWallet(r.Context(), pReq.CustomerXID)

	_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, c.Transformer.TransformInitWallet(token)))
}

func (c WalletControllerImpl) EnableWallet(w http.ResponseWriter, r *http.Request) {
	res := helper.PlugResponse(w)
	authPayload := r.Context().Value(constant.AuthorizationPayloadKey).(*token.Payload)

	wallet, err := c.Service.ActivateWallet(r.Context(), authPayload.CustomerXID)
	if err != nil {
		_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, helper.ErrData{Error: err.Error()}))
		return
	}

	_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, c.Transformer.TransformWallet(wallet)))
}

func (c WalletControllerImpl) DisableWallet(w http.ResponseWriter, r *http.Request) {
	res := helper.PlugResponse(w)
	authPayload := r.Context().Value(constant.AuthorizationPayloadKey).(*token.Payload)

	wallet, err := c.Service.DeactivateWallet(r.Context(), authPayload.CustomerXID)
	if err != nil {
		_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, helper.ErrData{Error: err.Error()}))
		return
	}

	_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, c.Transformer.TransformWalletDisable(wallet)))
}

func (c WalletControllerImpl) GetWallet(w http.ResponseWriter, r *http.Request) {
	res := helper.PlugResponse(w)
	authPayload := r.Context().Value(constant.AuthorizationPayloadKey).(*token.Payload)

	wallet, err := c.Service.GetWallet(r.Context(), authPayload.CustomerXID)
	if err != nil {
		_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, helper.ErrData{Error: err.Error()}))
		return
	}

	_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, c.Transformer.TransformWallet(wallet)))
}
