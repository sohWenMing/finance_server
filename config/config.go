package config

import (
	"database/sql"
	"os"

	envvars "github.com/sohWenMing/finance_server/env_vars"
	"github.com/sohWenMing/finance_server/internal/database/sqlc_generated"
)

/*
	the config holds is the main holder of all information of the application

	this includes:
	* Qeuries which are sql generated
*/

type Config struct {
	Queries   *sqlc_generated.Queries
	JwtSecret []byte
}

func (c *Config) RegisterJwtSecret(envPath string) error {
	envVarErr := envvars.LoadEnv(envPath)
	if envVarErr != nil {
		return envVarErr
	}
	secret := os.Getenv("JWT_SECRET")
	c.JwtSecret = []byte(secret)
	return nil
}

func (c *Config) RegisterQueries(db *sql.DB) {
	// at this point, the database should already be loaded, so we should be passing the db type into this function
	c.Queries = sqlc_generated.New(db)
}
