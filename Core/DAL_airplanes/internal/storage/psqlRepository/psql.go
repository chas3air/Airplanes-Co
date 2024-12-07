package psql

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PsqlStorage struct {
	DB *sql.DB
}
