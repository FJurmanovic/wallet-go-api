package filter

import "wallet-api/pkg/model"

type ApiFilter struct {
	model.Params
}

func NewApiFilter(params model.Params) *ApiFilter {
	return &ApiFilter{
		Params: params,
	}
}
