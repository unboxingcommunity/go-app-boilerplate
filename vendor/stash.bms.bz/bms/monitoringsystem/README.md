# Monitoring System

Go Agent allows you to monitor your Go applications currently supported with New Relic and Elastic APM. It helps you track transactions, outbound requests, database calls, and other parts of your. Go application's behavior and provides a running overview of garbage collection, goroutine activity, and memory use.

## Installation

To install run the following command

```sh
go get stash.bms.bz/bms/monitoringsystem
```

## NewRelic

The New Relic Go Agent allows you to monitor your Go applications with New
Relic. It helps you track transactions, outbound requests, database calls, and
other parts of your Go application's behavior and provides a running overview of
garbage collection, goroutine activity, and memory use.

**Required Below environment**

- NEWRELIC_APM
- NEWRELIC_KEY

[Example](_example/newrelic/main.go)

## Elastic

Elastic [APM](https://www.elastic.co/guide/en/apm/agent/go/current/configuration.html#config-service-node-name) is an application performance monitoring system built on the Elastic Stack. [ref][sre doc](https://confluence.bms.bz/display/AP/Elastic+APM)

**Required Below environment**

- ELASTIC_APM_SERVICE_NAME
- ELASTIC_APM_SERVER_URL AND ELASTIC_APM_SERVER_URLS
- ELASTIC_APM_SECRET_TOKEN: If your server requires a secret token for authentication, you must also set APM_TOKEN.

## gRCP Tracing Example

How to trace gRPC with APM (go tracer).

### Instrumentation

**Client Side**

```go
    apm, err := monitoringsystem.New(monitoringsystem.Elastic, true)
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		":15001",
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			apmgrpc.UnaryClientInterceptor(apmgrpc.WithAPM(apm)),
		),
		grpc.WithBlock(),
	)
```

You can see the details here: [\_client/main.go](_example/apmgrpc/client/main.go)

**Server Side**

```go
    apm, _ := monitoringsystem.New(monitoringsystem.Elastic, true)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			apmgrpc.UnaryServerInterceptor(apmgrpc.WithAPM(apm)),
		),
	)
```

You can see the details here: [\_server/main.go](_example/apmgrpc/main.go)

[Example](_example/elastic/main.go)

## TODO

- [ ] Need to do the step by step documentation of APM integration
