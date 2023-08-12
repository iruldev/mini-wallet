package service

const (
	DEPOSIT    = "deposit"
	WITHDRAWAL = "withdrawal"
)

type TransactionReq struct {
	Amount      string `json:"amount" validate:"required,numeric"`
	ReferenceID string `json:"reference_id" validate:"required"`
	Type        string
}
