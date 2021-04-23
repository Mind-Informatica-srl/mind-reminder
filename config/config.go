package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type EnvironmentProps struct {
	Db *gorm.DB
}

var Env EnvironmentProps

func init() {
	//istanzio il pool di connessioni
	db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=lamicolor dbname=lamicolor password=L4m1c0l0r sslmode=disable"), &gorm.Config{
		// Logger: appLogger,
	}) //gorm.Open("postgres", "host="+dbHost+" port=5432 user=lamicolor dbname=lamicolor password=L4m1c0l0r sslmode=disable")
	if err != nil {
		panic(err)
	}

	Env = EnvironmentProps{
		Db: db,
	}
}
