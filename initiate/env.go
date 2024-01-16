package initiate

import (
	"go-boilerplate-api/config"
	"os"

	"github.com/ralstan-vaz/go-errors"
)

// Version ...
var Version string

// Env Gets the enviroment variable from TIER
// By default the env will be development
func Env() (string, error) {
	env := os.Getenv("TIER")
	if env == "" {
		err := os.Setenv("TIER", config.ENVDevelopment)
		if err != nil {
			return "", errors.NewInternalError(err).SetCode("INITIATE.ENV.SETENV_FAILED")
		}
		env = os.Getenv("TIER")
	}

	return env, nil
}

// SetApmEnv ... sets various env flags for the apm
func SetApmEnv(env string) {
	os.Setenv("ELASTIC_APM_CAPTURE_HEADERS", "false")
}
