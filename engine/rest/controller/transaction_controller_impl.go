package controller

import (
	"net/http"

	"github.com/iruldev/mini-wallet/src/model/entity"

	"github.com/go-playground/validator/v10"
	"github.com/iruldev/mini-wallet/engine/rest/transformer"
	"github.com/iruldev/mini-wallet/src/constant"
	"github.com/iruldev/mini-wallet/src/helper"
	"github.com/iruldev/mini-wallet/src/service"
	"github.com/iruldev/mini-wallet/src/token"
)

type TransactionControllerImpl struct {
	Validator     *validator.Validate
	Service       service.TransactionService
	Transformer   transformer.TransactionTransformer
	WalletService service.WalletService
}

func NewTransactionController(
	validator *validator.Validate,
	service service.TransactionService,
	transformer transformer.TransactionTransformer,
	walletService service.WalletService,
) TransactionController {
	return &TransactionControllerImpl{
		Validator:     validator,
		Service:       service,
		Transformer:   transformer,
		WalletService: walletService,
	}
}

func (c TransactionControllerImpl) GetTransactions(w http.ResponseWriter, r *http.Request) {
	res := helper.PlugResponse(w)

	authPayload := r.Context().Value(constant.AuthorizationPayloadKey).(*token.Payload)

	wallet, err := c.WalletService.GetWallet(r.Context(), authPayload.CustomerXID)
	if err != nil {
		_ = res.ReplyCustom(http.StatusNotFound, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
		return
	}

	trns, err := c.Service.GetTransactions(r.Context(), wallet.CustomerXID)
	if err != nil || trns == nil {
		_ = res.ReplyCustom(http.StatusNotFound, helper.NewResponse(constant.FAILED, helper.ErrData{Error: "no transactions found"}))
		return
	}

	_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, c.Transformer.TransformerTransactions(trns)))

}

func (c TransactionControllerImpl) Deposit(w http.ResponseWriter, r *http.Request) {
	req := helper.PlugRequest(r, w)
	res := helper.PlugResponse(w)

	authPayload := r.Context().Value(constant.AuthorizationPayloadKey).(*token.Payload)

	pReq, _ := helper.ParseTo[service.TransactionReq](req)
	err := c.Validator.Struct(pReq)
	if err != nil {
		errF := helper.GetErrMsgField(err)
		_ = res.ReplyCustom(http.StatusBadRequest, helper.NewResponse(constant.FAILED, helper.ErrData{Error: errF}))
		return
	}

	pReq.Type = entity.DEPOSIT

	wallet, err := c.WalletService.GetWallet(r.Context(), authPayload.CustomerXID)
	if err != nil {
		_ = res.ReplyCustom(http.StatusNotFound, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
		return
	}

	trns, err := c.Service.Transaction(r.Context(), wallet.CustomerXID, pReq)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
		return
	}

	_ = res.ReplyCustom(http.StatusCreated, helper.NewResponse(constant.SUCCESS, c.Transformer.TransformerDeposit(trns)))
}

func (c TransactionControllerImpl) Withdrawal(w http.ResponseWriter, r *http.Request) {
	req := helper.PlugRequest(r, w)
	res := helper.PlugResponse(w)

	authPayload := r.Context().Value(constant.AuthorizationPayloadKey).(*token.Payload)
	pReq, _ := helper.ParseTo[service.TransactionReq](req)
	err := c.Validator.Struct(pReq)
	if err != nil {
		errF := helper.GetErrMsgField(err)
		_ = res.ReplyCustom(http.StatusBadRequest, helper.NewResponse(constant.FAILED, helper.ErrData{Error: errF}))
		return
	}

	pReq.Type = entity.WITHDRAWAL

	wallet, err := c.WalletService.GetWallet(r.Context(), authPayload.CustomerXID)
	if err != nil {
		_ = res.ReplyCustom(http.StatusNotFound, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
		return
	}

	trns, err := c.Service.Transaction(r.Context(), wallet.CustomerXID, pReq)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
		return
	}

	_ = res.ReplyCustom(http.StatusCreated, helper.NewResponse(constant.SUCCESS, c.Transformer.TransformerWithdrawal(trns)))
}
