package monitoringsystem

import (
	"net/http"
	"os"
	"sync"
)

const (
	envAppName     = "NEWRELIC_APP"
	envNewrelicKey = "NEWRELIC_KEY"
)

// agentType is used to identify APM (Application Performance Monitoring) proccess supported by BMS.
type agentType string

// Application monitoring supported by monitoring library
const (
	NewRelic agentType = "NewRelic"
	Elastic  agentType = "Elastic"
)

// Operation sets options for database tracing.
type Operation struct {
	// Instance holds the database instance name.
	Instance string
	// Statement holds the statement executed in the span,
	// e.g. "SELECT * FROM foo".
	Statement string
	// Type holds the database type, e.g. "sql".
	Type string
	// User holds the username used for database access.
	User string
}

// handler contanis all Application Performance Monitoring regarding the request context
type handler interface {
	StartTransaction(name string) (transaction interface{}, err error)
	StartWebTransaction(name string, rw http.ResponseWriter, req *http.Request) (transaction interface{}, err error)
	EndTransaction(transaction interface{}) (err error)
	StartSegment(segmentName string, Transaction interface{}) (segment interface{}, err error)
	EndSegment(segment interface{}) (err error)
	StartDataStoreSegment(segmentName string, transaction interface{}, operation string, collectionName string, operations ...Operation) (segment interface{}, err error)
	EndDataStoreSegment(segment interface{}) (err error)
	StartExternalSegment(transaction interface{}, URL string) (externalSegment interface{}, err error)
	StartExternalWebSegment(transaction interface{}, req *http.Request) (externalSegment interface{}, err error)
	EndExternalSegment(segment interface{}) error
	NoticeError(transaction interface{}, err error) error
	AddAttribute(transaction interface{}, key string, val interface{}) error
}

// Agent holds all application monitoring
type Agent struct {
	// application name is used by New Relic to link data across servers.
	appName string
	// monitorMode controls whether the agent will communicate with the APM servers or not
	// Setting this to be false is useful in testing and staging situations.
	enabled bool
	monitor handler
	mutex   *sync.Mutex
}

// Option are functions to set configuration for APM
// Following environment variables will be passed value used by default
// ELASTIC_APM_SERVICE_NAME
// ELASTIC_APM_SERVER_URL AND ELASTIC_APM_SERVER_URLS
// ELASTIC_APM_SECRET_TOKEN
// https://confluence.bms.bz/display/AP/Elastic+APM
//
// Recommend to use default which supported by SRE
type Option struct {
	ElasticServiceName  string
	ElasticServerAMPUrl string
	ElasticSecretKey    string
}

// New creates an Application and spawns goroutines to manage the
// aggregation and harvesting of data.  On success, a non-nil Application and a
// nil error are returned. On failure, a nil Application and a non-nil error
// are returned. Applications do not share global state, therefore it is safe
// to create multiple applications.
func New(monitoringType agentType, enableMonitoring bool, options ...Option) (*Agent, error) {

	if !enableMonitoring {
		return &Agent{appName: os.Getenv("CONFIG_APP"), enabled: enableMonitoring, mutex: &sync.Mutex{}}, nil
	}

	agent := &Agent{enabled: enableMonitoring, mutex: &sync.Mutex{}}
	switch monitoringType {
	case NewRelic:
		appName, licenseKey := os.Getenv(envAppName), os.Getenv(envNewrelicKey)
		if appName == "" {
			return nil, ErrInvalidAppName
		}
		if licenseKey == "" {
			return nil, ErrInvalidLicenseName
		}
		agent.appName = appName
		monitor, err := newRelicMonitor(agent.appName, licenseKey)
		if err != nil {
			return nil, err
		}
		agent.monitor = monitor
	case Elastic:
		// Following environment variables will be passed to application
		// https://confluence.bms.bz/display/AP/Elastic+APM
		name := os.Getenv("ELASTIC_APM_SERVICE_NAME")
		if name == "" {
			name = os.Getenv("CONFIG_APP")
		}
		opt := Option{}
		for _, option := range options {
			opt = option
		}
		if opt.ElasticServiceName != "" {
			name = opt.ElasticServiceName
		}
		if opt.ElasticServerAMPUrl == "" {
			// ELASTIC_APM_SERVER_URLS if specified, otherwise parses ELASTIC_APM_SERVER_URL if specified. If neither are specified, then the default localhost URL is returned.
			opt.ElasticServerAMPUrl = os.Getenv("ELASTIC_APM_SERVER_URL")
		}
		if opt.ElasticSecretKey == "" {
			opt.ElasticSecretKey = os.Getenv("ELASTIC_APM_SECRET_TOKEN")
		}
		monitor, _ := newElastic(name, opt.ElasticSecretKey, opt.ElasticServerAMPUrl)
		agent.monitor = monitor
	default:
		return nil, ErrUnsupported
	}
	return agent, nil
}

// Enable method allows if monitoring is disable data will not be pushed to APM
func (agent *Agent) Enable(monitoring bool) {
	if agent != nil || agent.monitor != nil {
		agent.mutex.Lock()
		agent.enabled = monitoring
		agent.mutex.Unlock()
	}
}

// isMonitoringEnabled determines whether the APM need enabled or not
// if monitoring is disable data will not be pushed to APM
func (agent *Agent) isMonitoringEnabled() bool {
	if agent == nil || agent.monitor == nil {
		return false
	}
	return agent.enabled
}

// StartTransaction returns a new Transaction with the specified
// name and type, and with the start time set to the current time.
// This is equivalent to calling StartTransactionOptions with a
// zero TransactionOptions.
func (agent *Agent) StartTransaction(transactionName string) (interface{}, error) {
	if !agent.isMonitoringEnabled() {
		return nil, nil
	}
	return agent.monitor.StartTransaction(transactionName)
}

// StartWebTransaction begins a web transaction.
// * The Transaction is considered a web transaction if an http.Request
//   is provided.
// * The transaction returned implements the http.ResponseWriter
//   interface.  Provide your ResponseWriter as a parameter and
//   then use the Transaction in its place to instrument the response
//   code and response headers.
func (agent *Agent) StartWebTransaction(transactionName string, rw http.ResponseWriter, req *http.Request) (interface{}, error) {
	if !agent.isMonitoringEnabled() {
		return nil, nil
	}
	return agent.monitor.StartWebTransaction(transactionName, rw, req)
}

// EndTransaction function finishes the transaction, stopping all further instrumentation.
//
// Calling End will set trans TransactionData field to nil, so callers
// must ensure tx is not updated after End returns.
func (agent *Agent) EndTransaction(trans interface{}, err error) error {
	if !agent.isMonitoringEnabled() {
		return nil
	}
	if err != nil {
		agent.monitor.NoticeError(trans, err)
	}
	return agent.monitor.EndTransaction(trans)
}

// StartSegment function instrument segments to the particular transaction.
//
// Input Parameters
//    Transaction(interface{})	: Transaction object which returned from StartTransaction function.
//    SegmentName(string)	: The name of the segment
//
// Output Parameters
//
//	interface{}		: segment object will be used as inputs for EndSegment function
// 	error			: error will be nil if function doesn't return errors else it will returns error.
func (agent *Agent) StartSegment(transObj interface{}, segmentName string) (interface{}, error) {
	if !agent.isMonitoringEnabled() {
		return nil, nil
	}
	return agent.monitor.StartSegment(segmentName, transObj)
}

// EndSegment function finishes the Segment.
//
// Segment(interface{}):Segment object which returned from StartSegment function.
func (agent *Agent) EndSegment(segment interface{}) error {
	if !agent.isMonitoringEnabled() {
		return nil
	}
	return agent.monitor.EndSegment(segment)
}

//StartDataStoreSegment function  is used to instrument calls to databases and object stores.
//
//Input Parameters
//    Transaction(interface{})	: Transaction object which returned from StartTransaction function.
//    SegmentName(string)	: The name of the segment
//    Operation(string)		: Operation is the relevant action, e.g. "SELECT" or "GET".
//    CollectionName(string)	: CollectionName is the table name or group name.
//
// Output Parameters
//
//	interface{}		: segment object will be used as inputs for EndDataStoreSegment function
// 	error			: error will be nil if function doesn't return errors else it will returns error.
func (agent *Agent) StartDataStoreSegment(transaction interface{}, segmentName string, operation string, collectionName string, operations ...Operation) (datastoreSegment interface{}, err error) {
	if !agent.isMonitoringEnabled() {
		return nil, nil
	}
	return agent.monitor.StartDataStoreSegment(segmentName, transaction, operation, collectionName, operations...)
}

// EndDataStoreSegment function finishes the datastore segment.
//
// Input Parameters
//    Segment(interface{}): Segment object which returned from StartDataStoreSegment function.
func (agent *Agent) EndDataStoreSegment(datastoreSegment interface{}) (err error) {
	if !agent.isMonitoringEnabled() {
		return nil
	}

	return agent.monitor.EndDataStoreSegment(datastoreSegment)
}

// StartExternalSegment function is used to instrument external calls.StartExternalSegment is recommended when you do not have access to an http.Request.
//
// Input Parameters
//    Transaction(interface{})	: Transaction object which returned from StartTransaction function.
//    URL(string)		: URL field should be used to indicate the endpoint.
//
// Output Parameters
//
//	interface{}		: segment object will be used as inputs for EndExternalSegment function
// 	error			: error will be nil if function doesn't return errors else it will returns error.
func (agent *Agent) StartExternalSegment(transObj interface{}, URL string) (externalSegment interface{}, err error) {
	if !agent.isMonitoringEnabled() {
		return nil, nil
	}
	return agent.monitor.StartExternalSegment(transObj, URL)
}

// StartExternalWebSegment function is used to instrument external calls.StartExternalWebSegment is recommended when you have access to an http.Request.
//
// Input Parameters
//    Transaction(interface{})	: Transaction object which returned from StartTransaction function.
//    req(*http.Request)		: URL field should be used to indicate the endpoint.
//
// Output Parameters
//
//	interface{}		: segment object will be used as inputs for EndExternalSegment function
// 	error			: error will be nil if function doesn't return errors else it will returns error.
func (agent *Agent) StartExternalWebSegment(transObj interface{}, req *http.Request) (externalSegment interface{}, err error) {
	if !agent.isMonitoringEnabled() {
		return nil, nil
	}

	return agent.monitor.StartExternalWebSegment(transObj, req)
}

// EndExternalSegment function finishes the external segment.
//
// Input Parameters
//    Segment(interface{}):Segment object which returned from StartExternalSegment function.
func (agent *Agent) EndExternalSegment(externalSegment interface{}) error {
	if !agent.isMonitoringEnabled() {
		return nil
	}

	return agent.monitor.EndExternalSegment(externalSegment)
}

// NoticeError function records an error.
//
// Input Parameters
// 	 transObj(interface{})	:	Transaction object which returned from StartTransaction function.
//	 err(error) 		: 	Error which need to be recorded.
func (agent *Agent) NoticeError(transObj interface{}, err error) error {
	if !agent.isMonitoringEnabled() {
		return nil
	}
	return agent.monitor.NoticeError(transObj, err)
}

// AddAttribute function adds a key value pair information to the current transaction.
//
// Input Parameters
//  transObj(interface{}) Transaction object which returned from StartTransaction function.
//	key(string)	key information specifies the name of the information attribute holds.
//	val(interface{}) This attribute specifies the value of the information attribute holds.
func (agent *Agent) AddAttribute(trans interface{}, key string, val interface{}) error {
	if !agent.isMonitoringEnabled() {
		return nil
	}
	return agent.monitor.AddAttribute(trans, key, val)
}
