package transformer

import "github.com/iruldev/mini-wallet/src/model/entity"

type InitWalletResponseData struct {
	Token string `json:"token"`
}

type WalletResponseData struct {
	Wallet WalletData `json:"wallet"`
}

type WalletDisableResponseData struct {
	Wallet WalletDisableData `json:"wallet"`
}

type WalletData struct {
	ID        string  `json:"id"`
	OwnedBy   string  `json:"owned_by"`
	Status    string  `json:"status"`
	EnabledAt string  `json:"enabled_at"`
	Balance   float64 `json:"balance"`
}

type WalletDisableData struct {
	ID         string  `json:"id"`
	OwnedBy    string  `json:"owned_by"`
	Status     string  `json:"status"`
	DisabledAt string  `json:"disabled_at"`
	Balance    float64 `json:"balance"`
}

type WalletTransformer interface {
	TransformInitWallet(token string) *InitWalletResponseData
	TransformWallet(data *entity.Wallet) *WalletResponseData
	TransformWalletDisable(data *entity.Wallet) *WalletDisableResponseData
}
