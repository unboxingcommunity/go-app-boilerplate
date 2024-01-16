package db

import (
	"go-boilerplate-api/config"
)

// Instances ... contains the interface layer of the different dbs
type Instances struct {
	MyDB MyDBInterface
}

// NewInstance creates an instance of initialized DBInstances
func NewInstance(conf config.IConfig) (*Instances, error) {
	dbInstances := &Instances{}

	myDBInstance, err := initMyDB(conf)
	if err != nil {
		return nil, err
	}

	// Sets db instance
	dbInstances.MyDB = myDBInstance

	return dbInstances, nil
}

// Simulates the initialization of a db connection
func initMyDB(config config.IConfig) (MyDBInterface, error) {
	return NewMyDB(), nil
}
