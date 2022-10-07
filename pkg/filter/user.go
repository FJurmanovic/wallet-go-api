package filter

import "wallet-api/pkg/model"

type UserFilter struct {
	model.Params
	BaseFilter
}

func NewUserFilter(params model.Params) *UserFilter {
	return &UserFilter{
		Params: params,
	}
}
