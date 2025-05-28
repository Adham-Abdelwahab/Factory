package util

import log "github.com/sirupsen/logrus"

type ResourceDetails struct {
	Name string
}

type DatabaseInterface interface {
	SetupDatabase() error
	GetResourceDetails(resource string) *ResourceDetails
}

func NewDatabase() (*DatabaseInterface, error) {
	var db DatabaseInterface = &db{}

	err := db.SetupDatabase()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &db, nil
}

type db struct{}

func (*db) SetupDatabase() error {
	return nil
}

func (*db) GetResourceDetails(resource string) *ResourceDetails {
	return &ResourceDetails{}
}
