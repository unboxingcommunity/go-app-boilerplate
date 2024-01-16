package ccms

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	ccmspb "stash.bms.bz/merchandise/utils/config/ccms/configpb"
)

// GetKey ...
func (c *GRPCClientLayer) GetKey(key string) (string, error) {
	key = strings.TrimSpace(key)
	value, err := c.Cli.Get(c.Ctx, &ccmspb.Input{Key: key})
	if err != nil {
		return "", err
	}
	output := value.Value
	return output, nil
}

// GetKey ...
func (c *LocalDataStore) GetKey(key string) (string, error) {
	if val, ok := c.GlobalKeyDSS[key]; ok {
		return val.(string), nil
	}
	return "", errors.New("Key Not Found")
}

func connectGRPC(host string, port string) (c *GRPCClientLayer, err error) {
	address := host + ":" + port
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c = &GRPCClientLayer{}
	c.Cli = ccmspb.NewCCMSClient(conn)
	c.Ctx = context.Background()
	return c, nil
}

// Init ...
func Init(host, port string, keysLocalStore map[string]interface{}, ccmsType string) (Provider, error) {
	if ccmsType == "" {
		ccmsType = getCCMSType()
	}
	if host == "" {
		host = getCCMSGRPCHOST()
	}
	if port == "" {
		port = getCCMSGRPCPORT()
	}

	provider, err := New(host, port, keysLocalStore, ccmsType)
	if err != nil {
		return nil, err
	}
	GProvider = provider
	return provider, nil
}

// New ...
func New(host, port string, keysLocalStore map[string]interface{}, ccmsType string) (Provider, error) {

	// Init Function can be overridden on Main GO
	//By default sidecar will be used regardless if error is not caught or handled
	switch ccmsType {
	case MAPKEYSCONFIG:
		{
			var datastore = new(LocalDataStore)
			datastore.GlobalKeyDSS = keysLocalStore
			return datastore, nil
		}
	default:
		{
			client, err := connectGRPC(host, port)
			if err != nil {
				return nil, err
			}
			return client, nil
		}
	}
}
