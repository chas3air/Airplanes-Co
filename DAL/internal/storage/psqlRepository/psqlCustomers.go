package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/chas3air/Airplanes-Co/DAL/internal/config"
	"github.com/chas3air/Airplanes-Co/DAL/internal/models"
	_ "github.com/lib/pq"
)

type PsqlCustomersStorage struct {
	PsqlStorage
}

func MustNewPsqlCustomersStorage(db *sql.DB) PsqlCustomersStorage {
	const op = "DAL.internal.storage.psqlRepository.newCustomersStorage"
	err := db.Ping()
	if err != nil {
		log.Println("Table of customers is unavailable: " + err.Error())
		panic(fmt.Errorf("%s: %w", op, err))
	}

	return PsqlCustomersStorage{
		PsqlStorage: NewPsqlStorage(db),
	}
}

func (s PsqlCustomersStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.GetAll"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of customers is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of customers is unavailable, file: %s: %w", op, err)
	}

	rows, err := s.DB.QueryContext(ctx, `SELECT * FROM `+config.PSQL_CUSTOMERS_TABLE_NAME+`;`)
	if err != nil {
		log.Println("Error querying customers:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()
	customers := make([]models.Customer, 0, 10)
	customer := models.Customer{}

	for rows.Next() {
		err := rows.Scan(&customer.Id, &customer.Login,
			&customer.Password, &customer.Role, &customer.Surname, &customer.Name)
		if err != nil {
			log.Println("Error scanning row:", err.Error())
			continue
		}
		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		log.Println("Row error:", err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Count of retrieved customers: %v\n", len(customers))
	return customers, nil
}

func (s PsqlCustomersStorage) GetById(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.GetById"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of customers is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of customers is unavailable, file: %s: %w", op, err)
	}

	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+config.PSQL_CUSTOMERS_TABLE_NAME+`
		WHERE id = $1;
		`, id)

	customer := models.Customer{}

	err = row.Scan(&customer.Id, &customer.Login,
		&customer.Password, &customer.Role, &customer.Surname, &customer.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No customer found with id=%d\n", id)
			return nil, fmt.Errorf("%s: No flight found with id=%d", op, id)
		}
		log.Println("Error scanning customer:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrived customer`s id=%d\n", customer.Id)
	return customer, nil
}

func (s PsqlCustomersStorage) GetByLoginAndPassword(ctx context.Context, login string, password string) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.GetByLoginAndPassword"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of customers is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of customers is unavailable, file: %s: %w", op, err)
	}

	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+config.PSQL_CUSTOMERS_TABLE_NAME+`
		WHERE login = $1 AND password = $2
		`, login, password)

	customer := models.Customer{}

	err = row.Scan(&customer.Id, &customer.Login,
		&customer.Password, &customer.Role, &customer.Surname, &customer.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No customer found with login=%s and password=%s\n", login, password)
			return nil, fmt.Errorf("%s: запись с login=%s и password=%s не найдена", op, login, password)
		}
		log.Println("Error scanning flight:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrived customer`s id=%d\n", customer.Id)
	return customer, nil
}

func (s PsqlCustomersStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.Insert"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of customers is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of customers is unavailable, file: %s: %w", op, err)
	}

	customer := innerObj.(models.Customer)
	var id int = -1

	err = s.DB.QueryRowContext(ctx,
		`INSERT INTO `+config.PSQL_CUSTOMERS_TABLE_NAME+`
		(login, password, role, surname, name) VALUES ($1, $2, $3, $4, $5) RETURNING id
		`, customer.Login, customer.Password, customer.Role, customer.Surname, customer.Name).Scan(&id)

	if err != nil {
		log.Println("Error inserting customer:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return s.GetById(ctx, id)
}

func (s PsqlCustomersStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.Update"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of customers is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of customers is unavailable, file: %s: %w", op, err)
	}

	customer := innerObj.(models.Customer)

	_, err = s.DB.ExecContext(ctx,
		`UPDATE `+config.PSQL_CUSTOMERS_TABLE_NAME+`
		SET login = $1, password = $2, role = $3, surname = $4, name = $5
		WHERE id = $6;
	`, customer.Login, customer.Password, customer.Role, customer.Surname, customer.Name, customer.Id)

	if err != nil {
		log.Println("Error updating customer:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Updated customer with id=%d\n", customer.Id)
	return s.GetById(ctx, customer.Id)
}

func (s PsqlCustomersStorage) Delete(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.Delete"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of customers is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of customers is unavailable, file: %s: %w", op, err)
	}

	customer, err := s.GetById(ctx, id)
	if err != nil {
		log.Println("Error getting flight by id:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+config.PSQL_CUSTOMERS_TABLE_NAME+`
		WHERE id = $1;
		`, id)

	if err != nil {
		log.Println("Error deleting flight:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Deleted flight with id=%d\n", id)
	return customer, nil
}
