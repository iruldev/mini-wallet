package service

type TransactionReq struct {
	Amount      string `json:"amount" validate:"required,numeric"`
	ReferenceID string `json:"reference_id" validate:"required"`
	Type        string
}
