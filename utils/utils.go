package utils

import (
	"log"

	"github.com/ztrue/tracerr"
)

func CheckError(err error) {
	if err != nil {
		tracerr.PrintSourceColor(err)
		log.Fatal(err)
	}
}
