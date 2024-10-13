package psql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chas3air/Airplanes-Co/DAL/internal/config"
)

type PsqlStorage struct {
	DB *sql.DB
}

func NewPsqlStorage(db *sql.DB) PsqlStorage {
	return PsqlStorage{DB: db}
}

func InitDB() *sql.DB {
	const op = "DAL.internal.storage.psqlRepository.InitDB"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.PSQL_DB_USER,
		config.PSQL_DB_PASSWORD,
		config.PSQL_DB_HOST,
		config.PSQL_DB_PORT,
		config.PSQL_DB_DBNAME)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(fmt.Errorf("%s: %s", op, err))
	}

	if err := db.Ping(); err != nil {
		log.Fatalln(fmt.Errorf("%s: %s", op, err))
	}

	return db
}
