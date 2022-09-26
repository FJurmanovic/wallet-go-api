package filter

import "wallet-api/pkg/model"

type TransactionFilter struct {
	model.Params
	BaseFilter
	NoPending           bool
	TransactionStatusId string
}

func NewTransactionFilter(params model.Params) *TransactionFilter {
	return &TransactionFilter{
		Params: params,
	}
}
