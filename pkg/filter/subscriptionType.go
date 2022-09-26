package filter

import "wallet-api/pkg/model"

type SubscriptionTypeFilter struct {
	model.Params
	BaseFilter
}

func NewSubscriptionTypeFilter(params model.Params) *SubscriptionTypeFilter {
	return &SubscriptionTypeFilter{
		Params: params,
	}
}
