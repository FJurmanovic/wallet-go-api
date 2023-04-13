package service

import (
	"context"
	"log"
	"wallet-api/pkg/model"
	"wallet-api/pkg/repository"
	"wallet-api/pkg/utl/common"
)

type CurrencyService struct {
	repository *repository.CurrencyRepository
	logger     *log.Logger
}

func NewCurrencyService(repository *repository.CurrencyRepository, logger *log.Logger) *CurrencyService {
	return &CurrencyService{
		repository: repository,
		logger:     logger,
	}
}

/*
GetFirst

Gets first row from Currency table.

	   	Args:
	   		context.Context: Application context
		Returns:
			model.CurrencyModel: Currency object from database.
*/
func (as CurrencyService) Sync(ctx context.Context) {
	resp, err := common.Fetch[model.ExchangeBody]("GET", "https://api.exchangerate-api.com/v4/latest/euro")
	if err != nil {
		as.logger.Println(err)
		return
	}
	m := resp.Rates.(map[string]interface{})

	rates := new([]model.Rate)

	for k, v := range m {
		rate := new(model.Rate)
		rate.Code = k
		rate.Rate = v.(float64)
		*rates = append(*rates, *rate)
	}
	as.repository.SyncBulk(ctx, rates)
}
