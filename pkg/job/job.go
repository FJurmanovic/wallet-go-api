package job

import (
	"log"
	"os"
	"wallet-api/pkg/service"
	"wallet-api/pkg/utl/common"

	"go.uber.org/dig"
)

/*
InitializeJobs

Initializes Dependency Injection modules and registers Jobs

	Args:
		*dig.Container: Dig Container
*/
func InitializeJobs(c *dig.Container) {
	file, err := os.OpenFile("job.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	common.CheckError(err)
	logger := log.New(file, "Job: ", log.Ldate|log.Ltime|log.Lshortfile)

	jobContainer := c.Scope("job")
	jobContainer.Provide(func() *log.Logger {
		return logger
	})

	service.InitializeServices(jobContainer)

	jobContainer.Invoke(NewCurrencyJob)
}
