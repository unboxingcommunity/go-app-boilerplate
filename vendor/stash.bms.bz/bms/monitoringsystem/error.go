package monitoringsystem

import "errors"

// application monitoring error's
var (
	ErrInvalidAppName     = errors.New("Application name is not set")
	ErrInvalidLicenseName = errors.New("APM License is not set")
	ErrUnsupported        = errors.New("Unsupported monitoring system type")
	ErrNotInitialized     = errors.New("APM application object is nil")
	ErrInvalidTrans       = errors.New("APM transaction object is nil")
	ErrInvalidRequest     = errors.New("HTTP request object didn't provided")
	ErrInvalidSegment     = errors.New("APM segment object is nil")
	ErrInvalidDataType    = errors.New("Invalid data type")
)
