package filter

import "wallet-api/pkg/model"

type WalletHeaderFilter struct {
	model.Params
	BaseFilter
}

func NewWalletHeaderFilter(params model.Params) *WalletHeaderFilter {
	return &WalletHeaderFilter{
		Params: params,
	}
}
