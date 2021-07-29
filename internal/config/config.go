// Package configuration contain the current configuration
// in use by the project
package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
func Create(dsn string, production bool) (*Data, error) {
	var err error
	config := &Data{}

	gormConfig := &gorm.Config{}
	if production {
		gormConfig.Logger = logger.Discard
	} else {
		gormConfig.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Enable color
			},
		)

	}

	if config.DB, err = gorm.Open(postgres.Open(dsn), gormConfig); err != nil {
		return nil, err
	}

	return config, err
}
