package filter

import "wallet-api/pkg/model"

type UserFilter struct {
	model.Params
	BaseFilter
	Email string
}

func NewUserFilter(params model.Params) *UserFilter {
	return &UserFilter{
		Params: params,
	}
}
