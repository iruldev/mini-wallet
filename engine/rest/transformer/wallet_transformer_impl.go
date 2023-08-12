package transformer

import (
	"time"

	"github.com/iruldev/mini-wallet/src/model/entity"
)

type WalletTransformerImpl struct {
}

func NewWalletTransformer() WalletTransformer {
	return &WalletTransformerImpl{}
}

func (WalletTransformerImpl) TransformInitWallet(token string) *InitWalletResponseData {
	return &InitWalletResponseData{Token: token}
}

func (WalletTransformerImpl) TransformWallet(data *entity.Wallet) *WalletResponseData {
	enabledAt := time.Now().Format(time.RFC3339)
	if data.ActivatedAt != nil {
		enabledAt = data.ActivatedAt.Format(time.RFC3339)
	}
	return &WalletResponseData{WalletData{
		ID:        data.ID.String(),
		OwnedBy:   data.CustomerXID,
		Status:    "enabled",
		EnabledAt: enabledAt,
		Balance:   data.Balance.InexactFloat64(),
	}}
}

func (WalletTransformerImpl) TransformWalletDisable(data *entity.Wallet) *WalletDisableResponseData {
	return &WalletDisableResponseData{WalletDisableData{
		ID:         data.ID.String(),
		OwnedBy:    data.CustomerXID,
		Status:     "disabled",
		DisabledAt: time.Now().Format(time.RFC3339),
		Balance:    data.Balance.InexactFloat64(),
	}}
}
