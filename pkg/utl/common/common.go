package common

import (
	"log"
	"time"
	"wallet-api/pkg/models"

	"github.com/google/uuid"
)

func CheckError(err error) {
	if err != nil {
		log.Fatalf("Error occured. %v", err)
	}
}

func CreateDbModel() models.CommonModel {
	date := time.Now()
	dbModel := models.CommonModel{
		Id:          uuid.NewString(),
		DateCreated: date,
		DateUpdated: date,
	}
	return dbModel
}
