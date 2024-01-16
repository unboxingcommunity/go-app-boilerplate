// Package apmgrpc provides interceptors for tracing monitoring gRPC.
package apmgrpc

import (
	"context"
	"fmt"
	"go-boilerplate-api/apm"
	log "go-boilerplate-api/pkg/utils/logger"

	"google.golang.org/grpc"
	"stash.bms.bz/bms/monitoringsystem"
)

var (
	defaultOptions = &options{}
)

// options options for creating a request context object
type options struct {
	apm *monitoringsystem.Agent
}

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

// Option sets options for server-side tracing.
type Option func(*options)

// WithAPM customizes the function for monitoring the request performance
func WithAPM(apm *monitoringsystem.Agent) Option {
	return func(o *options) {
		o.apm = apm
	}
}

// UnaryServerInterceptor returns a grpc.UnaryServerInterceptor that
// traces gRPC requests with the given options.
//
// The interceptor will trace transactions with the "grpc" type for each
// incoming request.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOptions(opts)
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if o.apm != nil {
			// Starts an APM transaction
			tx, err := o.apm.StartTransaction(info.FullMethod)
			if err != nil || tx == nil {
				log.Error("GO-BOILERPLATE.GRPC.APM_TRANS_INIT_FAIL", "Transaction failed", log.Priority1, nil, map[string]interface{}{"error": err.Error()})
			}

			// Stores transaction details in context
			ctx = context.WithValue(ctx, apm.TransactionKey, tx)

			// Ends transaction
			defer func(ctx context.Context, opts *options) {
				if err != nil {
					err = fmt.Errorf("panic: %v", err)
				}
				opts.apm.EndTransaction(tx, err)
			}(ctx, o)
		}

		resp, err = handler(ctx, req)
		return resp, err

	}
}
