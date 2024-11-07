package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	psqlUser      = os.Getenv("PSQL_DB_USER")
	psqlPassword  = os.Getenv("PSQL_DB_PASSWORD")
	psqlHost      = os.Getenv("PSQL_DB_HOST")
	psqlPort      = os.Getenv("PSQL_DB_PORT")
	psqlDBName    = os.Getenv("PSQL_DB_DBNAME")
	psqlTableName = os.Getenv("PSQL_TABLE_NAME")
)

type PsqlStorage struct {
	DB *sql.DB
}

func NewPsqlStorage(db *sql.DB) PsqlStorage {
	log.Println("Creating new PostgreSQL storage")
	return PsqlStorage{DB: db}
}

// InitDB initializes the database connection.
func InitDB() *sql.DB {
	const op = "DAL.internal.storage.psqlRepository.InitDB"
	log.Println("Initializing database connection")

	// Проверка переменных окружения
	if psqlUser == "" || psqlPassword == "" || psqlHost == "" || psqlPort == "" || psqlDBName == "" || psqlTableName == "" {
		log.Panic("Database configuration variables are not set")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		psqlUser,
		psqlPassword,
		psqlHost,
		psqlPort,
		psqlDBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(fmt.Errorf("%s: failed to open connection to database: %s", op, err))
	}

	log.Println("Pinging the database to check connection")
	for i := 0; i < 5; i++ {
		if err := db.Ping(); err == nil {
			break
		}
		log.Println("Waiting for database to be ready...")
		time.Sleep(200 * time.Millisecond)
	}

	log.Println("Database connection established successfully")
	return db
}
