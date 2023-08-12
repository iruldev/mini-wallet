package transformer

import (
	"time"

	"github.com/iruldev/mini-wallet/src/model/entity"
)

type TransactionTransformerImpl struct {
}

func NewTransactionTransformer() TransactionTransformer {
	return &TransactionTransformerImpl{}
}

func (TransactionTransformerImpl) TransformerDeposit(data *entity.Transaction) *DepositTransactionResponseData {
	return &DepositTransactionResponseData{DepositTransactionData{
		ID:          data.ID.String(),
		DepositedBy: data.CustomerXID,
		Status:      data.Status,
		DepositedAt: data.At.Format(time.RFC3339),
		Amount:      data.Amount.InexactFloat64(),
		ReferenceID: data.ReferenceID,
	}}
}

func (TransactionTransformerImpl) TransformerWithdrawal(data *entity.Transaction) *WithdrawalTransactionResponseData {
	return &WithdrawalTransactionResponseData{WithdrawalTransactionData{
		ID:          data.ID.String(),
		WithdrawnBy: data.CustomerXID,
		Status:      data.Status,
		WithdrawnAt: data.At.Format(time.RFC3339),
		Amount:      data.Amount.InexactFloat64(),
		ReferenceID: data.ReferenceID,
	}}
}
