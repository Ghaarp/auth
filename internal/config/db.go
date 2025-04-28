package config

import (
	"fmt"
	"log"
	"os"
)

const (
	dbEnvName         = "PG_DATABASE_NAME"
	dbUserEnvName     = "PG_USER"
	dbPasswordEnvName = "PG_PASSWORD"
	dbHostEnvName     = "PG_HOST"
	dbPortEnvName     = "PG_PORT"
)

type DBConfig interface {
	DSN() string
}

type dbConfig struct {
	dbName     string
	dbUserName string
	dbPassword string
	dbHost     string
	dbPort     string
}

func (dbconf *dbConfig) DSN() string {
	res := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", dbconf.dbHost, dbconf.dbPort, dbconf.dbName, dbconf.dbUserName, dbconf.dbPassword)
	log.Print(res)
	return res
}

func NewDBConfig() (DBConfig, error) {
	dbName := os.Getenv(dbEnvName)
	if len(dbName) == 0 {
		return nil, fmt.Errorf("db name not found in .env")
	}

	dbUserName := os.Getenv(dbUserEnvName)
	if len(dbUserName) == 0 {
		return nil, fmt.Errorf("db username not found in .env")
	}

	dbPassword := os.Getenv(dbPasswordEnvName)
	if len(dbPassword) == 0 {
		return nil, fmt.Errorf("db password not found in .env")
	}

	dbHost := os.Getenv(dbHostEnvName)
	if len(dbHost) == 0 {
		return nil, fmt.Errorf("db port not found in .env")
	}

	dbPort := os.Getenv(dbPortEnvName)
	if len(dbPort) == 0 {
		return nil, fmt.Errorf("db port not found in .env")
	}

	return &dbConfig{
		dbName:     dbName,
		dbUserName: dbUserName,
		dbPassword: dbPassword,
		dbHost:     dbHost,
		dbPort:     dbPort,
	}, nil
}
