package transformer

import "github.com/iruldev/mini-wallet/src/model/entity"

type DepositTransactionResponseData struct {
	Deposit DepositTransactionData `json:"deposit"`
}

type WithdrawalTransactionResponseData struct {
	Withdrawal WithdrawalTransactionData `json:"withdrawal"`
}

type DepositTransactionData struct {
	ID          string  `json:"id"`
	DepositedBy string  `json:"deposited_by"`
	Status      string  `json:"status"`
	DepositedAt string  `json:"deposited_at"`
	Amount      float64 `json:"amount"`
	ReferenceID string  `json:"reference_id"`
}

type WithdrawalTransactionData struct {
	ID          string  `json:"id"`
	WithdrawnBy string  `json:"withdrawn_by"`
	Status      string  `json:"status"`
	WithdrawnAt string  `json:"withdrawn_at"`
	Amount      float64 `json:"amount"`
	ReferenceID string  `json:"reference_id"`
}

type TransactionTransformer interface {
	TransformerDeposit(data *entity.Transaction) *DepositTransactionResponseData
	TransformerWithdrawal(data *entity.Transaction) *WithdrawalTransactionResponseData
}
