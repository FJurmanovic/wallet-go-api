package common

import "log"

func CheckError(err error) {
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
}
