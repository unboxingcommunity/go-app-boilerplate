package main

import (
	"go-boilerplate-api/initiate"
	log "go-boilerplate-api/pkg/utils/logger"
	"go-boilerplate-api/shared"

	"github.com/ralstan-vaz/go-errors"
)

// Version ...
var Version string

func main() {
	// Sets the version flag
	shared.VERSION = Version

	// Initializes logger
	log.InitLogger()

	// Initialize the app
	err := initiate.Initialize()
	if err != nil {
		newErr := errors.Get(err)
		// If an error is encountered in main its critical and the app exits
		log.Fatal(newErr.Code, newErr.Description, newErr.Source)
	}
}
