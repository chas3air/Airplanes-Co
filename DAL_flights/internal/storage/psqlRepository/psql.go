package psql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chas3air/Airplanes-Co/DAL_flights/internal/config"
	_ "github.com/lib/pq"
)

type PsqlStorage struct {
	DB *sql.DB
}

func NewPsqlStorage(db *sql.DB) PsqlStorage {
	log.Println("Creating new PostgreSQL storage")
	return PsqlStorage{DB: db}
}

func InitDB() *sql.DB {
	const op = "DAL.internal.storage.psqlRepository.InitDB"
	log.Println("Initializing database connection")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.PSQL_DB_USER,
		config.PSQL_DB_PASSWORD,
		config.PSQL_DB_HOST,
		config.PSQL_DB_PORT,
		config.PSQL_DB_DBNAME)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(fmt.Errorf("%s: failed to open connection to database: %s", op, err))
	}

	log.Println("Pinging the database to check connection")
	if err := db.Ping(); err != nil {
		log.Fatalln(fmt.Errorf("%s: failed to ping database: %s", op, err))
	}

	log.Println("Database connection established successfully")
	return db
}
