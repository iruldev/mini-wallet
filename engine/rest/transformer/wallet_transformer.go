package transformer

type InitWalletResponseData struct {
	Token string `json:"token"`
}

type WalletTransformer interface {
	TransformInitWallet(token string) *InitWalletResponseData
}
