package config

import (
	log "go-boilerplate-api/pkg/utils/logger"
	"io/ioutil"

	"gopkg.in/yaml.v3"
	"stash.bms.bz/merchandise/utils"
	"stash.bms.bz/merchandise/utils/config/ccms"
)

var configModel *Config

// NewConfig gets the configuration based on the environment passed
func NewConfig(env string) (IConfig, error) {

	configFile := "config/tier/" + env + ".yaml"
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if env != "development" && env != "sit" && env != "docker" && env != "testing" {
		cfg, err := InitCCMS(env)
		if err != nil {
			return nil, err
		}

		// anonymous function to fetch from ccms
		getCcmsValue := func(key string) (string, error) {
			value, err := cfg.GetKey(key)
			if err != nil {
				log.Fatal("CONFIG.KEY.NOT.FOUND", "Key Not Found", log.Priority1, nil, map[string]interface{}{key: err.Error()})
				return "", err
			}
			return value, err
		}

		// Binds config based on ccms and present values
		err = utils.BindConfig(bytes, &configModel, "ccms", getCcmsValue)
		if err != nil {
			return nil, err
		}

		return &IConfigModel{model: configModel}, nil
	}

	err = yaml.Unmarshal(bytes, &configModel)
	if err != nil {
		return nil, err
	}

	// Returns
	return &IConfigModel{model: configModel}, nil
}

// Get implements the interface function for IConfig
func (ic *IConfigModel) Get() *Config {
	return ic.model
}

// InitCCMS ...
func InitCCMS(env string) (ccms.Provider, error) {
	ccmsConfig, err := ccms.Init("", "", nil, "")
	if err != nil {
		return nil, err
	}
	return ccmsConfig, nil
}
