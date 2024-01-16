package utils

import (
	"log"

	errors "stash.bms.bz/merchandise/go-errors"
)

func init() {
	er := errors.Parse(GetErrors())
	if er != nil {
		log.Panicln(er)
	}
}
