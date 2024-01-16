package apm

import (
	"context"
	"net/http"

	"stash.bms.bz/bms/monitoringsystem"
)

// HandlerInterface ... A wrapper interface on top of the apm (monitoringsystem) to help out during testing
type HandlerInterface interface {
	StartTransaction(name string) (transaction interface{}, err error)
	EndTransaction(transaction interface{}, er error) (err error)
	StartSegment(ctx context.Context, segmentName string) (segment interface{}, err error)
	EndSegment(segment interface{}) (err error)
	StartDataStoreSegment(ctx context.Context, segmentName string, operation string, collectionName string, operations ...monitoringsystem.Operation) (datastoreSegment interface{}, err error)
	EndDataStoreSegment(segment interface{}) (err error)
	StartExternalSegment(ctx context.Context, URL string) (externalSegment interface{}, err error)
	EndExternalSegment(segment interface{}) error
	StartExternalWebSegment(ctx context.Context, req *http.Request) (externalSegment interface{}, err error)
	NoticeError(transaction interface{}, err error) error
	AddAttribute(transaction interface{}, key string, val interface{}) error
}

// Handler ...
type Handler struct {
}

// NewApmHandler ...
func NewApmHandler() HandlerInterface {
	return &Handler{}
}

// NOTE : We could return nil,nil from every function if APM == nil , then we don't need to mock , but I'm a bit split on this atm

// StartTransaction ... Interface function that calls the internal apm functions
func (a *Handler) StartTransaction(name string) (transaction interface{}, err error) {
	return APM.StartTransaction(name)
}

// EndTransaction .. Interface function that calls the internal apm functions
func (a *Handler) EndTransaction(transaction interface{}, er error) (err error) {
	return APM.EndTransaction(transaction, er)
}

// StartSegment .. Interface function that calls the internal apm functions
func (a *Handler) StartSegment(ctx context.Context, segmentName string) (segment interface{}, err error) {
	transaction := FromContext(ctx)
	return APM.StartSegment(transaction, segmentName)
}

// EndSegment .. Interface function that calls the internal apm functions
func (a *Handler) EndSegment(segment interface{}) (err error) {
	return APM.EndSegment(segment)
}

// StartDataStoreSegment .. Interface function that calls the internal apm functions
func (a *Handler) StartDataStoreSegment(ctx context.Context, segmentName string, operation string, collectionName string, operations ...monitoringsystem.Operation) (datastoreSegment interface{}, err error) {
	transaction := FromContext(ctx)
	return APM.StartDataStoreSegment(transaction, segmentName, operation, collectionName, operations...)
}

// EndDataStoreSegment .. Interface function that calls the internal apm functions
func (a *Handler) EndDataStoreSegment(segment interface{}) (err error) {
	return APM.EndDataStoreSegment(segment)
}

// StartExternalSegment .. Interface function that calls the internal apm functions
func (a *Handler) StartExternalSegment(ctx context.Context, URL string) (externalSegment interface{}, err error) {
	transaction := FromContext(ctx)
	return APM.StartExternalSegment(transaction, URL)
}

// EndExternalSegment .. Interface function that calls the internal apm functions
func (a *Handler) EndExternalSegment(segment interface{}) error {
	return APM.EndExternalSegment(segment)
}

// StartExternalWebSegment .. Interface function that calls the internal apm functions
func (a *Handler) StartExternalWebSegment(ctx context.Context, req *http.Request) (externalSegment interface{}, err error) {
	transaction := FromContext(ctx)
	return APM.StartExternalWebSegment(transaction, req)
}

// NoticeError .. Interface function that calls the internal apm functions
func (a *Handler) NoticeError(transaction interface{}, err error) error {
	return APM.NoticeError(transaction, err)
}

// AddAttribute .. Interface function that calls the internal apm functions
func (a *Handler) AddAttribute(transaction interface{}, key string, val interface{}) error {
	return APM.AddAttribute(transaction, key, val)
}
