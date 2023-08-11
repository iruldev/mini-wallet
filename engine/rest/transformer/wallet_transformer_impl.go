package transformer

type WalletTransformerImpl struct {
}

func NewWalletTransformer() WalletTransformer {
	return &WalletTransformerImpl{}
}

func (WalletTransformerImpl) TransformInitWallet(token string) *InitWalletResponseData {
	return &InitWalletResponseData{Token: token}
}
