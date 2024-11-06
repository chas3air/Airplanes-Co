package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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
		os.Getenv("PSQL_DB_USER"),
		os.Getenv("PSQL_DB_PASSWORD"),
		os.Getenv("PSQL_DB_HOST"),
		os.Getenv("PSQL_DB_PORT"),
		os.Getenv("PSQL_DB_DBNAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(fmt.Errorf("%s: failed to open connection to database: %s", op, err))
	}

	log.Println("Pinging the database to check connection")
	for i := 0; i < 5; i++ {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Println("Waiting for database to be ready...")
		time.Sleep(200 * time.Millisecond)
	}

	log.Println("Database connection established successfully")
	return db
}
