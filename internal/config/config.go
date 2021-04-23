package config

import (
	"os"

	mindlogger "github.com/Mind-Informatica-srl/mind-logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type EnvironmentProps struct {
	db *gorm.DB
}

func (e *EnvironmentProps) GetDb(logger *mindlogger.AppLogger) *gorm.DB {
	e.db.Logger = logger
	return e.db
}

var Env EnvironmentProps

func init() {
	dbHost := os.Getenv("DB_HOST")
	//istanzio il pool di connessioni
	db, err := gorm.Open(postgres.Open("host="+dbHost+" port=5432 user=lamicolor dbname=lamicolor password=L4m1c0l0r sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	Env = EnvironmentProps{
		db: db,
	}
}
