package utils

import (
	"github.com/mitchellh/mapstructure"
	"github.com/ralstan-vaz/go-errors"
)

// Bind ... uses mapstructure internally to bind input to output (does not use tags)
func Bind(input interface{}, output interface{}) error {
	err := mapstructure.Decode(input, output)
	if err != nil {
		return errors.NewInternalError(err).SetCode("PKG.UTILS.DECODE_ERROR")
	}

	return nil
}
