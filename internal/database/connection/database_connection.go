package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	envvars "github.com/sohWenMing/finance_server/env_vars"
)

func connect(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectToDB(envPath string) (*sql.DB, error) {
	envVarErr := envvars.LoadEnv(envPath)
	if envVarErr != nil {
		return nil, envVarErr
	}
	DbString := os.Getenv("DB_FINANCE")
	db, err := connect(DbString)
	if err != nil {
		return nil, err
	}
	return db, nil
}
