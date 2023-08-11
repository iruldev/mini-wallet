package service

type InitWalletReq struct {
	CustomerXID string `json:"customer_xid" validate:"required"`
}
