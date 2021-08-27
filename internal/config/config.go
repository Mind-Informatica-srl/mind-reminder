// Package config contain the current configuration
// in use by the project
package config

import (
	"fmt"

	"gorm.io/gorm"
)

// Data contain all the configuration parameters
type Data struct {
	// DB is the used connection pool
	DB *gorm.DB
}

// currentConfiguration is the default configuration
var currentConfiguration *Data

// ErrorAlreadySetup is raised when this project is started up
// more than one times
var ErrorAlreadySetup = fmt.Errorf("configuration already set")

// ErrorMissingConfiguration is raised when we try to access to currentConfiguration that is not configured
var ErrorMissingConfiguration = fmt.Errorf("missing configuration")

// Setup create the current configuration.
func Setup(config *Data) error {
	if currentConfiguration != nil {
		return ErrorAlreadySetup
	}
	currentConfiguration = config
	return nil
}

// Current gets the current configurfation
func Current() *Data {
	return currentConfiguration
}

// Create create a configuration
func Create(db *gorm.DB) (*Data, error) {
	currentConfiguration = &Data{
		DB: db,
	}

	return currentConfiguration, nil
}
