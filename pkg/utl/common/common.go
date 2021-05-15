package common

import (
	"log"
)

func CheckError(err error) {
	if err != nil {
		log.Fatalf("Error occured. %v", err)
	}
}
