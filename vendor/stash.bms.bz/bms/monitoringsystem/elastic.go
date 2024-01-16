package monitoringsystem

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.elastic.co/apm"
	"go.elastic.co/apm/transport"
	"google.golang.org/grpc/status"
)

// elastic holds all the tracer manages the sampling and sending of transactions to
// Elastic APM.
type elastic struct {
	tracer *apm.Tracer
}

func newElastic(name, token string, urls ...string) (*elastic, error) {
	defer panicRecover()
	transport, err := transport.NewHTTPTransport()
	if err != nil {
		return nil, err
	}
	if len(urls) != 0 {
		for _, u := range urls {
			u = strings.TrimSpace(u)
			if u == "" {
				continue
			}
			if url, err := url.Parse(u); err == nil {
				transport.SetServerURL(url)
			}
		}
	}
	if token != "" {
		transport.SetSecretToken(token)
	}
	tracer := apm.DefaultTracer
	tracer.Transport = transport
	tracer.Service.Name = name
	tracer.Service.Environment = os.Getenv("TIER")
	return &elastic{
		tracer: tracer,
	}, nil
}

func (n *elastic) isValid() error {
	if n == nil || n.tracer == nil {
		return ErrNotInitialized
	}
	return nil
}

func (n *elastic) StartTransaction(name string) (transaction interface{}, err error) {
	defer panicRecover()
	if err := n.isValid(); err != nil {
		return nil, err
	}
	opts := apm.TransactionOptions{
		Start: time.Now(),
	}
	txn := n.tracer.StartTransactionOptions(name, "request", opts)
	txn.Result = "Success"
	return txn, nil
}

func (n *elastic) StartWebTransaction(name string, w http.ResponseWriter, req *http.Request) (transaction interface{}, err error) {
	defer panicRecover()
	if err := n.isValid(); err != nil {
		return nil, err
	}
	opts := apm.TransactionOptions{
		Start: time.Now(),
	}
	txn := n.tracer.StartTransactionOptions(name, req.Method, opts)
	// function relates to server-side requests. Various proxy
	// forwarding headers are taken into account to reconstruct the URL,
	// and determining the client address.
	if req != nil && req.URL != nil {
		txn.Context.SetHTTPRequest(req)
	}
	txn.Result = "Success"
	return txn, nil
}

func (n *elastic) EndTransaction(transaction interface{}) error {
	defer panicRecover()
	if err := n.isValid(); err != nil {
		return err
	}
	if transaction == nil {
		return ErrInvalidTrans
	}
	if txn, ok := transaction.(*apm.Transaction); ok && txn != nil {
		txn.End()
		return nil
	}
	return ErrInvalidDataType
}

func (n *elastic) StartSegment(name string, transaction interface{}) (interface{}, error) {
	defer panicRecover()
	if transaction == nil {
		return nil, ErrInvalidTrans
	}
	if txn, ok := transaction.(*apm.Transaction); ok && txn != nil {
		opts := apm.SpanOptions{
			Start: time.Now(),
		}
		span := txn.StartSpanOptions(name, "", opts)
		return span, nil
	}
	return nil, ErrInvalidDataType
}

func (n *elastic) EndSegment(segment interface{}) error {
	defer panicRecover()
	if err := n.isValid(); err != nil {
		return err
	}
	if segment == nil {
		return ErrInvalidSegment
	}
	if span, ok := segment.(*apm.Span); ok && span != nil {
		span.End()
		return nil
	}
	return ErrInvalidDataType
}

func (n *elastic) StartDataStoreSegment(name string, transaction interface{}, operation string, collectionName string, operations ...Operation) (interface{}, error) {
	defer panicRecover()
	if err := n.isValid(); err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, ErrInvalidTrans
	}
	if txn, ok := transaction.(*apm.Transaction); ok && txn != nil {
		opt := Operation{}
		for _, o := range operations {
			opt = o
		}
		opts := apm.SpanOptions{
			Start: time.Now(),
		}
		span := txn.StartSpanOptions(name, fmt.Sprintf("db.%s.%s", collectionName, operation), opts)
		if !span.Dropped() {
			span.Context.SetDatabase(apm.DatabaseSpanContext{
				Instance:  opt.Instance,
				Statement: opt.Statement,
				Type:      name,
				User:      opt.User,
			})
		}
		return span, nil
	}
	return nil, ErrInvalidDataType
}

func (n *elastic) EndDataStoreSegment(segment interface{}) error {
	defer panicRecover()
	return n.EndSegment(segment)
}

func (n *elastic) StartExternalSegment(transaction interface{}, url string) (interface{}, error) {
	defer panicRecover()
	if err := n.isValid(); err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, ErrInvalidTrans
	}
	if txn, ok := transaction.(*apm.Transaction); ok && txn != nil {
		opts := apm.SpanOptions{
			Start: time.Now(),
		}
		span := txn.StartSpanOptions(url, "", opts)
		return span, nil
	}
	return nil, ErrInvalidDataType
}

func (n *elastic) StartExternalWebSegment(transaction interface{}, req *http.Request) (interface{}, error) {
	defer panicRecover()
	if err := n.isValid(); err != nil {
		return nil, err
	}
	if transaction == nil {
		return nil, ErrInvalidTrans
	}
	if txn, ok := transaction.(*apm.Transaction); ok && txn != nil {
		opts := apm.SpanOptions{
			Start: time.Now(),
		}
		if req != nil && req.URL != nil {
			span := txn.StartSpanOptions(req.URL.Path, "", opts)
			return span, nil
		}
		return nil, ErrInvalidRequest
	}
	return nil, ErrInvalidDataType
}

func (n *elastic) EndExternalSegment(segment interface{}) error {
	defer panicRecover()
	return n.EndSegment(segment)
}

func (n *elastic) NoticeError(transaction interface{}, err error) error {
	defer panicRecover()
	if err := n.isValid(); err != nil {
		return err
	}
	if transaction == nil {
		return ErrInvalidTrans
	}
	if txn, ok := transaction.(*apm.Transaction); ok && txn != nil {
		txn.Result = "Failed"
		if s, ok := status.FromError(err); ok {
			txn.Result = s.Code().String()
		}
		if err != nil {
			n.tracer.NewError(err).Send()
		}
		return nil
	}
	return ErrInvalidDataType
}

func (n *elastic) AddAttribute(transaction interface{}, key string, value interface{}) error {
	defer panicRecover()
	if err := n.isValid(); err != nil {
		return err
	}
	if transaction == nil {
		return ErrInvalidTrans
	}
	if txn, ok := transaction.(*apm.Transaction); ok && txn != nil {
		txn.Context.SetLabel(key, value)
		return nil
	}
	return ErrInvalidDataType
}
