package filter

import "wallet-api/pkg/model"

type WalletFilter struct {
	model.Params
	BaseFilter
}

func NewWalletFilter(params model.Params) *WalletFilter {
	return &WalletFilter{
		Params: params,
	}
}
