package utils

import (
	"log"
)

func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
