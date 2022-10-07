package filter

import "wallet-api/pkg/model"

type TransactionStatusFilter struct {
	model.Params
	BaseFilter
	Status string
}

func NewTransactionStatusFilter(params model.Params) *TransactionStatusFilter {
	return &TransactionStatusFilter{
		Params: params,
	}
}
