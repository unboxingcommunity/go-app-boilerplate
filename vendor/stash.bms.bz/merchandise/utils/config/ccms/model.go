package ccms

import (
	"context"

	ccmspb "stash.bms.bz/merchandise/utils/config/ccms/configpb"
)

const (
	SIDECAR       string = "sidecar"
	CLIENT        string = "client"
	MAPKEYSCONFIG string = "mapstruct"
)

// GProvider ...
var GProvider Provider

// GRPCClientLayer ...
type GRPCClientLayer struct {
	Cli ccmspb.CCMSClient
	Ctx context.Context
}

// LocalDataStore ...
type LocalDataStore struct {
	GlobalKeyDSS map[string]interface{}
}

// Provider ...
type Provider interface {
	GetKey(string) (string, error)
}
