package grpc

import (
	ierror "errors"
	"go-boilerplate-api/apis/grpc/utils"
	"go-boilerplate-api/apis/middleware/apmgrpc"
	"go-boilerplate-api/apm"
	log "go-boilerplate-api/pkg/utils/logger"
	"go-boilerplate-api/shared"
	"net"
	"sync"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/ralstan-vaz/go-errors"
	"google.golang.org/grpc"
)

const (
	// grpcNetwork : The network must always belong to ["tcp", "tcp4", "tcp6", "unix" or "unixpacket"].
	grpcNetwork string = "tcp"
)

// StartServer starts the grpc server using the dependencies passed to it
func StartServer(deps *shared.Deps, wg *sync.WaitGroup, fatalError chan error) error {
	address := deps.Config.Get().Server.GRPC.Address

	var server *grpc.Server
	// Add required opts
	recoveryOpts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(handlePanic),
	}

	apmOpts := []apmgrpc.Option{
		apmgrpc.WithAPM(apm.APM),
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			apmgrpc.UnaryServerInterceptor(apmOpts...),
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		)),
	}

	// Creates new GRPC server
	server = grpc.NewServer(opts...)

	registerService(server, deps)

	// Creates TCP listener at a particular port
	lis, err := net.Listen(grpcNetwork, address)
	if err != nil {
		newErr := errors.NewInternalError(err).SetCode("APIS.GRPC.LISTENER_FAILED")
		fatalError <- newErr
	}

	// Logs server address
	log.Debug("GRPC Server listening on : " + address + ", Version: " + shared.VERSION)

	// Links server to the listener
	err = server.Serve(lis)
	if err != nil {
		newErr := errors.NewInternalError(err).SetCode("APIS.GRPC.LISTENER_LINK_FAILED")
		fatalError <- newErr
	}

	// Go routine finished (can also be deferred)
	wg.Done()

	return nil
}

// handlePanic handles unhandled panics by sending an error response for GRPC handlers
func handlePanic(p interface{}) error {
	err, ok := p.(error)
	if !ok {
		newErr := ierror.New("Panic recovery failed to parse error : " + err.Error())
		return utils.HandleError(&newErr)
	}

	return utils.HandleError(&err)
}
