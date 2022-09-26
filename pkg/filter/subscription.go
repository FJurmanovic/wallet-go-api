package filter

import "wallet-api/pkg/model"

type SubscriptionFilter struct {
	model.Params
	BaseFilter
}

func NewSubscriptionFilter(params model.Params) *SubscriptionFilter {
	return &SubscriptionFilter{
		Params: params,
	}
}
