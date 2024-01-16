package apm

import (
	"context"
	"go-boilerplate-api/config"
	log "go-boilerplate-api/pkg/utils/logger"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ralstan-vaz/go-errors"
	"stash.bms.bz/bms/monitoringsystem"
)

// APM - instance of APM monitoring system
var APM *monitoringsystem.Agent

const (
	// TransactionKey -key which stores the apm transaction in context
	TransactionKey string = "ApmTransactionKey"
)

// Initialize ... initializes the apm instance
func Initialize() {
	var app *monitoringsystem.Agent
	var enableMonitoring = true

	env := os.Getenv("TIER")
	if env == config.ENVDevelopment || env == config.ENVDocker {
		enableMonitoring = false
	}

	app, err := monitoringsystem.New(monitoringsystem.Elastic, enableMonitoring, monitoringsystem.Option{
		ElasticServiceName:  os.Getenv("ELASTIC_APM_SERVICE_NAME"),
		ElasticServerAMPUrl: os.Getenv("ELASTIC_APM_SERVER_URL"),
	})

	if err != nil {
		// turn off monitoring if there is error , but do not crash the app
		if app == nil {
			app, _ = monitoringsystem.New(monitoringsystem.Elastic, false)
		}
		newErr := errors.NewInternalError(err)
		log.Error(newErr.Code, "APM.INITIALIZE_FAILED", log.Priority1, newErr.Source)
		return
	}

	APM = app
}

// FromContext - used to fetch transaction key from the context which keeps track of APM transaction
// Input
//		ctx - is either the regular context or gin context
// Output
//		transactionKey - returns transaction key(or nil) which keeps track of APM transaction
func FromContext(ctx context.Context) interface{} {
	ct, _ := ctx.(*gin.Context)
	if ct != nil {
		vl, _ := ct.Get(TransactionKey)
		return vl
	}

	return ctx.Value(TransactionKey)
}
