package config

import (
	"fmt"
	"os"
)

const (
	DEFAULT_LIMIT_TIME = 5
)

var (
	psqlUser      = os.Getenv("PSQL_DB_USER")
	psqlPassword  = os.Getenv("PSQL_DB_PASSWORD")
	psqlHost      = os.Getenv("PSQL_DB_HOST")
	psqlPort      = os.Getenv("PSQL_DB_PORT")
	psqlDBName    = os.Getenv("PSQL_DB_DBNAME")
	psqlTableName = os.Getenv("PSQL_TABLE_NAME")
)

var ConnStr = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
	psqlUser,
	psqlPassword,
	psqlHost,
	psqlPort,
	psqlDBName,
)
