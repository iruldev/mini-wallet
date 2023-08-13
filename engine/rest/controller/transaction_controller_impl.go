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
	Validator   *validator.Validate
	Service     service.TransactionService
	Transformer transformer.TransactionTransformer
}

func NewTransactionController(
	validator *validator.Validate,
	service service.TransactionService,
	transformer transformer.TransactionTransformer,
) TransactionController {
	return &TransactionControllerImpl{
		Validator:   validator,
		Service:     service,
		Transformer: transformer,
	}
}

func (c TransactionControllerImpl) GetTransactions(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (c TransactionControllerImpl) Deposit(w http.ResponseWriter, r *http.Request) {
	req := helper.PlugRequest(r, w)
	res := helper.PlugResponse(w)

	authPayload := r.Context().Value(constant.AuthorizationPayloadKey).(*token.Payload)

	pReq, _ := helper.ParseTo[service.TransactionReq](req)
	pReq.Type = entity.DEPOSIT
	//err := c.Validator.Struct(pReq)
	//fmt.Println("pReq", pReq)
	//fmt.Println("err", err)
	trns, err := c.Service.Transaction(r.Context(), authPayload.CustomerXID, pReq)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
		return
	}

	_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, c.Transformer.TransformerDeposit(trns)))
}

func (c TransactionControllerImpl) Withdrawal(w http.ResponseWriter, r *http.Request) {
	req := helper.PlugRequest(r, w)
	res := helper.PlugResponse(w)

	authPayload := r.Context().Value(constant.AuthorizationPayloadKey).(*token.Payload)
	pReq, _ := helper.ParseTo[service.TransactionReq](req)
	pReq.Type = entity.WITHDRAWAL

	trns, err := c.Service.Transaction(r.Context(), authPayload.CustomerXID, pReq)
	if err != nil {
		_ = res.ReplyCustom(http.StatusBadRequest, helper.NewResponse(constant.FAILED, helper.ErrData{Error: err.Error()}))
		return
	}

	_ = res.ReplyCustom(http.StatusOK, helper.NewResponse(constant.SUCCESS, c.Transformer.TransformerWithdrawal(trns)))
}
