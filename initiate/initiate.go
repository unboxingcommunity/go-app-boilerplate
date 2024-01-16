package initiate

import (
	"go-boilerplate-api/apis"
	"go-boilerplate-api/apm"
	"go-boilerplate-api/config"
	"go-boilerplate-api/pkg/clients/db"
	grpcPkg "go-boilerplate-api/pkg/clients/grpc"
	httpPkg "go-boilerplate-api/pkg/clients/http"
	log "go-boilerplate-api/pkg/utils/logger"
	"go-boilerplate-api/shared"
)

// Initialize will initialize all the dependencies and the servers.
// Dependencies include config, external connections(grpc, http)
func Initialize() error {
	// Initializes logger
	log.InitLogger()

	env, err := Env()
	if err != nil {
		return err
	}

	log.Info("Enviroment: " + env)

	// Sets apm env
	SetApmEnv(env)

	// Initializes APM based on environment
	if env != config.ENVDevelopment && env != config.ENVDocker && env != config.ENVSit {
		// remove sit when building an actual app
		apm.Initialize()
	}

	// Gets config
	conf, err := config.NewConfig(env)
	if err != nil {
		return err
	}

	// Initializes the DB connections
	dbInstances, err := db.NewInstance(conf)
	if err != nil {
		return err
	}

	// Initializes the GRPC connections
	grpcCons, err := grpcPkg.NewConnections(conf)
	if err != nil {
		return err
	}

	// Initializes apm Handler
	handler := apm.NewApmHandler()

	// loads all common dependencies
	dependencies := shared.Deps{
		Config:        conf,
		Database:      dbInstances,
		GrpcConn:      grpcCons,
		HTTPRequester: httpPkg.NewRequest(),
		Apm:           handler,
	}

	// Initializes servers
	err = apis.InitServers(&dependencies)
	if err != nil {
		return err
	}

	// Returns
	return nil
}
