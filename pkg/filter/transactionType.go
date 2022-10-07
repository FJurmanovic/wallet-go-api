package filter

import "wallet-api/pkg/model"

type TransactionTypeFilter struct {
	model.Params
	BaseFilter
}

func NewTransactionTypeFilter(params model.Params) *TransactionTypeFilter {
	return &TransactionTypeFilter{
		Params: params,
	}
}
