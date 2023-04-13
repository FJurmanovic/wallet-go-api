package job

import (
	"context"
	"log"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"github.com/go-co-op/gocron"
)

type CurrencyController struct {
	service   *service.CurrencyService
	logger    *log.Logger
	scheduler *gocron.Scheduler
}

/*
NewCurrencyJob

Initializes CurrencyJob.

	Args:
		*services.CurrencyService: Currency service
		*gin.RouterGroup: Gin Router Group
	Returns:
		*CurrencyJob: Job for "Currency" route interactions
*/
func NewCurrencyJob(as *service.CurrencyService, scheduler *gocron.Scheduler, logger *log.Logger) *CurrencyController {
	currencyScheduler := scheduler.Tag("currency")

	wc := &CurrencyController{
		service:   as,
		logger:    logger,
		scheduler: currencyScheduler,
	}

	_, err := currencyScheduler.Every(1).Days().Do(wc.Sync)
	common.CheckError(err)
	currencyScheduler.StartAsync()

	log.Println("CurrencyJob started")

	return wc
}

func (wc *CurrencyController) Sync() {
	wc.logger.Println("CurrencyJob: Syncing currencies")
	ctx := context.Background()
	wc.service.Sync(ctx)
	wc.logger.Println("CurrencyJob: Syncing currencies done")
}
