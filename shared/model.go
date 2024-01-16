package shared

import (
	"go-boilerplate-api/apm"
	"go-boilerplate-api/config"
	"go-boilerplate-api/pkg/clients/db"
	grpcPkg "go-boilerplate-api/pkg/clients/grpc"
	httpPkg "go-boilerplate-api/pkg/clients/http"
)

// VERSION keeps the version no. (commit id) for global use
var VERSION string

// Deps ... is a shared dependencies struct that contains common singletons
type Deps struct {
	Config        config.IConfig
	Database      *db.Instances
	GrpcConn      grpcPkg.IGrpcConnections
	HTTPRequester httpPkg.IRequest
	Apm           apm.HandlerInterface
}
