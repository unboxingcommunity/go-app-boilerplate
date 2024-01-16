package monitoringsystem

import (
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

// newRelicHandler application monitoring
type newRelicHandler struct {
	// AppName is used by New Relic to link data across servers.
	//
	// https://docs.newrelic.com/docs/apm/new-relic-apm/installation-configuration/naming-your-application
	appName string
	// License is your New Relic license key.
	//
	// https://docs.newrelic.com/docs/accounts/install-new-relic/account-setup/license-key
	license     string
	application newrelic.Application
}

func newRelicMonitor(appName, license string) (*newRelicHandler, error) {
	nr := &newRelicHandler{appName: appName, license: license}
	config := newrelic.NewConfig(nr.appName, nr.license)
	application, err := newrelic.NewApplication(config)
	if err != nil {
		return nil, err
	}
	nr.application = application
	return nr, nil
}

func (n *newRelicHandler) isValid() error {
	if n == nil || n.application == nil {
		return ErrNotInitialized
	}
	return nil
}

func (n *newRelicHandler) StartTransaction(transactionName string) (transaction interface{}, err error) {
	if err = n.isValid(); err != nil {
		return
	}
	transaction = n.application.StartTransaction(transactionName, nil, new(http.Request))
	return
}

func (n *newRelicHandler) StartWebTransaction(transactionName string, rw http.ResponseWriter, req *http.Request) (transaction interface{}, err error) {
	if err = n.isValid(); err != nil {
		return
	}
	transaction = n.application.StartTransaction(transactionName, rw, req)
	return
}

func (n *newRelicHandler) EndTransaction(transaction interface{}) error {
	if err := n.isValid(); err != nil {
		return err
	}
	if transaction == nil {
		return ErrInvalidTrans
	}
	if txn, ok := transaction.(newrelic.Transaction); ok {
		return txn.End()
	}
	return ErrInvalidDataType
}

func (n *newRelicHandler) StartSegment(name string, transaction interface{}) (interface{}, error) {
	if err := n.isValid(); err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, ErrInvalidTrans
	}
	if txn, ok := transaction.(newrelic.Transaction); ok {
		return newrelic.Segment{
			Name:      name,
			StartTime: newrelic.StartSegmentNow(txn),
		}, nil
	}
	return nil, ErrInvalidDataType
}

func (n *newRelicHandler) EndSegment(segment interface{}) error {
	if err := n.isValid(); err != nil {
		return err
	}
	if segment == nil {
		return ErrInvalidSegment
	}
	if objSegment, ok := segment.(newrelic.Segment); ok {
		return objSegment.End()
	}
	return ErrInvalidDataType
}

func (n *newRelicHandler) StartDataStoreSegment(name string, transaction interface{}, operation string, collectionName string, operations ...Operation) (interface{}, error) {
	if err := n.isValid(); err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, ErrInvalidTrans
	}

	if txn, ok := transaction.(newrelic.Transaction); ok {
		return newrelic.DatastoreSegment{
			Product:    newrelic.DatastoreProduct(name),
			Collection: collectionName,
			Operation:  operation,
			StartTime:  newrelic.StartSegmentNow(txn),
		}, nil
	}

	return nil, ErrInvalidDataType
}

func (n *newRelicHandler) EndDataStoreSegment(segment interface{}) error {
	if err := n.isValid(); err != nil {
		return err
	}
	if segment == nil {
		return ErrInvalidSegment
	}
	if objSegment, ok := segment.(newrelic.DatastoreSegment); ok {
		return objSegment.End()
	}
	return ErrInvalidDataType
}

func (n *newRelicHandler) StartExternalSegment(transaction interface{}, url string) (interface{}, error) {
	if err := n.isValid(); err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, ErrInvalidTrans
	}
	if txn, ok := transaction.(newrelic.Transaction); ok {
		return newrelic.ExternalSegment{
			StartTime: newrelic.StartSegmentNow(txn),
			URL:       url,
		}, nil
	}
	return nil, ErrInvalidDataType
}

func (n *newRelicHandler) StartExternalWebSegment(transaction interface{}, req *http.Request) (interface{}, error) {
	if err := n.isValid(); err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, ErrInvalidTrans
	}
	if txn, ok := transaction.(newrelic.Transaction); ok {
		return newrelic.StartExternalSegment(txn, req), nil
	}
	return nil, ErrInvalidDataType
}

func (n *newRelicHandler) EndExternalSegment(segment interface{}) error {
	if err := n.isValid(); err != nil {
		return err
	}
	if segment == nil {
		return ErrInvalidSegment
	}
	if objSegment, ok := segment.(newrelic.ExternalSegment); ok {
		return objSegment.End()
	}
	return ErrInvalidDataType
}

func (n *newRelicHandler) NoticeError(transaction interface{}, err error) error {
	if err := n.isValid(); err != nil {
		return err
	}
	if transaction == nil {
		return ErrInvalidTrans
	}
	if txn, ok := transaction.(newrelic.Transaction); ok {
		return txn.NoticeError(err)
	}
	return ErrInvalidDataType
}

func (n *newRelicHandler) AddAttribute(transaction interface{}, key string, val interface{}) error {
	if err := n.isValid(); err != nil {
		return err
	}
	if transaction == nil {
		return ErrInvalidTrans
	}
	if txn, ok := transaction.(newrelic.Transaction); ok {
		return txn.AddAttribute(key, val)
	}
	return ErrInvalidDataType
}
