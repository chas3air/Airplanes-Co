package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/config"
	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/models"
)

type PsqlCustomersStorage struct {
	PsqlStorage
}

// MustNewPsqlCustomersStorage initializes a new PsqlCustomersStorage instance.
// It checks the database connection and panics if the customers table is unavailable.
func MustNewPsqlCustomersStorage(db *sql.DB) PsqlCustomersStorage {
	const op = "DAL.internal.storage.psqlRepository.newCustomersStorage"
	err := db.Ping()
	if err != nil {
		log.Println("Customers table is unavailable: " + err.Error())
		panic(fmt.Errorf("%s: %w", op, err))
	}

	return PsqlCustomersStorage{
		PsqlStorage: NewPsqlStorage(db),
	}
}

// GetAll retrieves all customers from the database.
// It returns a slice of Customer models and an error, if any occurs.
func (s PsqlCustomersStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.GetAll"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Customers table is unavailable: " + err.Error())
		return nil, fmt.Errorf("%s: %w", op, err)
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

	log.Printf("Retrieved %v customers\n", len(customers))
	return customers, nil
}

// GetById retrieves a customer by its ID.
// It returns the Customer model and an error if no customer is found or if an error occurs.
func (s PsqlCustomersStorage) GetById(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.GetById"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Customers table is unavailable")
		return nil, fmt.Errorf("customers table is unavailable, file: %s: %w", op, err)
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
			log.Printf("No customer found with ID=%d\n", id)
			return nil, fmt.Errorf("%s: no customer found with ID=%d", op, id)
		}
		log.Println("Error scanning customer:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrieved customer with ID=%d\n", customer.Id)
	return customer, nil
}

// GetByLoginAndPassword retrieves a customer by their login and password.
// It returns the Customer model and an error if no matching customer is found or if an error occurs.
func (s PsqlCustomersStorage) GetByLoginAndPassword(ctx context.Context, login string, password string) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.GetByLoginAndPassword"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Customers table is unavailable: " + err.Error())
		return nil, fmt.Errorf("customers table is unavailable, file: %s: %w", op, err)
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
			log.Printf("No customer found with login=%s and provided password\n", login)
			return nil, fmt.Errorf("%s: no customer found with login=%s", op, login)
		}
		log.Println("Error scanning customer:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrieved customer with ID=%d\n", customer.Id)
	return customer, nil
}

// Insert adds a new customer to the database.
// It returns the inserted Customer model and an error if the insertion fails.
func (s PsqlCustomersStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.Insert"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Customers table is unavailable: " + err.Error())
		return nil, fmt.Errorf("customers table is unavailable, file: %s: %w", op, err)
	}

	customer := innerObj.(models.Customer)
	var id int

	err = s.DB.QueryRowContext(ctx,
		`INSERT INTO `+config.PSQL_CUSTOMERS_TABLE_NAME+`
		(login, password, role, surname, name) VALUES ($1, $2, $3, $4, $5) RETURNING id
		`, customer.Login, customer.Password, customer.Role, customer.Surname, customer.Name).Scan(&id)

	if err != nil {
		log.Println("Error inserting customer:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	log.Printf("Inserted customer with ID=%d\n", id)
	return s.GetById(ctx, id)
}

// Update modifies an existing customer in the database.
// It returns the updated Customer model and an error if the update fails.
func (s PsqlCustomersStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.Update"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Customers table is unavailable: " + err.Error())
		return nil, fmt.Errorf("customers table is unavailable, file: %s: %w", op, err)
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

	log.Printf("Updated customer with ID=%d\n", customer.Id)
	return s.GetById(ctx, customer.Id)
}

// Delete removes a customer from the database by its ID.
// It returns the deleted Customer model and an error if the deletion fails.
func (s PsqlCustomersStorage) Delete(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlCustomers.Delete"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Customers table is unavailable: " + err.Error())
		return nil, fmt.Errorf("customers table is unavailable, file: %s: %w", op, err)
	}

	customer, err := s.GetById(ctx, id)
	if err != nil {
		log.Println("Error fetching customer by ID:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+config.PSQL_CUSTOMERS_TABLE_NAME+`
		WHERE id = $1;
		`, id)

	if err != nil {
		log.Println("Error deleting customer:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Deleted customer with ID=%d\n", id)
	return customer, nil
}
