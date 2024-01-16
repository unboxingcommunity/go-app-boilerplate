package utils

import (
	log "go-boilerplate-api/pkg/utils/logger"

	"github.com/ralstan-vaz/go-errors"
	"github.com/ralstan-vaz/go-errors/grpc"
	"google.golang.org/grpc/status"
)

// HandleError formats, logs and sets a GRPC response for the error
func HandleError(errObj *error) error {
	if *errObj == nil {
		return nil
	}

	err := errors.Get(*errObj)

	if err.Message == "" {
		err.Message = "Something Went Wrong"
	}

	log.Error(err.Code, err.Description, log.Priority1, err.Source)

	statusCode := grpc.StatusCode(err)

	return status.New(statusCode, err.Description).Err()
}
