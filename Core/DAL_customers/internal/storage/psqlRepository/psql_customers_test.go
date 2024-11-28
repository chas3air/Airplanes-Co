package psql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chas3air/Airplanes-Co/Core/DAL_customers/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TODO дописать тесты

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlCustomersStorage(db)

	t.Run("Success", func(t *testing.T) {
		expectedCustomers := []models.Customer{
			{Id: uuid.New(), Login: "aisachenkov", Password: "Bebrolov", Role: "admin", Surname: "Isachenkov", Name: "Artem"},
			{Id: uuid.New(), Login: "dlazin", Password: "loger123", Role: "user", Surname: "Lazin", Name: "Denis"},
		}

		rows := sqlmock.NewRows([]string{"id", "login", "password", "role", "surname", "name"}).
			AddRow(expectedCustomers[0].Id, expectedCustomers[0].Login, expectedCustomers[0].Password, expectedCustomers[0].Role, expectedCustomers[0].Surname, expectedCustomers[0].Name).
			AddRow(expectedCustomers[1].Id, expectedCustomers[1].Login, expectedCustomers[1].Password, expectedCustomers[1].Role, expectedCustomers[1].Surname, expectedCustomers[1].Name)

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + ";").
			WillReturnRows(rows)

		customers, err := storage.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCustomers, customers)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("QueryError", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + ";").
			WillReturnError(fmt.Errorf("query error"))

		customers, err := storage.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, customers)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectation: %s", err)
		}
	})

}

func TestGetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlCustomersStorage(db)

	t.Run("Success", func(t *testing.T) {
		customerID := uuid.New()
		expectedCustomer := models.Customer{
			Id:       customerID,
			Login:    "aisachenkov",
			Password: "Bebrolov",
			Role:     "admin",
			Surname:  "Isachenkov",
			Name:     "Artem",
		}

		row := sqlmock.NewRows([]string{"id", "login", "password", "role", "surname", "name"}).
			AddRow(expectedCustomer.Id, expectedCustomer.Login, expectedCustomer.Password, expectedCustomer.Role, expectedCustomer.Surname, expectedCustomer.Name)

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1").
			WithArgs(customerID).
			WillReturnRows(row)

		customer, err := storage.GetById(context.Background(), customerID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCustomer, customer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		customerID := uuid.New()

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1").
			WithArgs(customerID).
			WillReturnError(sql.ErrNoRows)

		customer, err := storage.GetById(context.Background(), customerID)
		assert.Error(t, err)
		assert.Nil(t, customer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		customerID := uuid.New()

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1").
			WithArgs(customerID).
			WillReturnRows(sqlmock.NewRows([]string{"id"})) // Only one column

		customer, err := storage.GetById(context.Background(), customerID)
		assert.Error(t, err)
		assert.Nil(t, customer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetByLoginAndPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlCustomersStorage(db)

	t.Run("Success", func(t *testing.T) {
		expectedCustomer := models.Customer{
			Id:       uuid.New(),
			Login:    "aisachenkov",
			Password: "Bebrolov",
			Role:     "admin",
			Surname:  "Isachenkov",
			Name:     "Artem",
		}

		row := sqlmock.NewRows([]string{"id", "login", "password", "role", "surname", "name"}).
			AddRow(expectedCustomer.Id, expectedCustomer.Login, expectedCustomer.Password, expectedCustomer.Role, expectedCustomer.Surname, expectedCustomer.Name)

		mock.ExpectQuery("SELECT \\* FROM "+os.Getenv("PSQL_TABLE_NAME")+" WHERE login = \\$1 AND password = \\$2").
			WithArgs(expectedCustomer.Login, expectedCustomer.Password).
			WillReturnRows(row)

		customer, err := storage.GetByLoginAndPassword(context.Background(), expectedCustomer.Login, expectedCustomer.Password)
		assert.NoError(t, err)
		assert.Equal(t, expectedCustomer, customer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM "+os.Getenv("PSQL_TABLE_NAME")+" WHERE login = \\$1 AND password = \\$2").
			WithArgs("invalid_login", "invalid_password").
			WillReturnError(sql.ErrNoRows)

		customer, err := storage.GetByLoginAndPassword(context.Background(), "invalid_login", "invalid_password")
		assert.Error(t, err)
		assert.Nil(t, customer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("ScanError", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM "+os.Getenv("PSQL_TABLE_NAME")+" WHERE login = \\$1 AND password = \\$2").
			WithArgs("login", "password").
			WillReturnRows(sqlmock.NewRows([]string{"id"})) // Only one column

		customer, err := storage.GetByLoginAndPassword(context.Background(), "login", "password")
		assert.Error(t, err)
		assert.Nil(t, customer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlCustomersStorage(db)

	t.Run("Success", func(t *testing.T) {
		customer := models.Customer{
			Login:    "aisachenkov",
			Password: "Bebrolov",
			Role:     "admin",
			Surname:  "Isachenkov",
			Name:     "Artem",
		}

		id := uuid.New()

		// Set up expectation for the insert operation
		mock.ExpectQuery("INSERT INTO "+os.Getenv("PSQL_TABLE_NAME")+
			" \\(login, password, role, surname, name\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id").
			WithArgs(customer.Login, customer.Password, customer.Role, customer.Surname, customer.Name).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

		// Set up expectation for retrieving the customer by ID after insertion
		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1").
			WithArgs(id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "role", "surname", "name"}).
				AddRow(id, customer.Login, customer.Password, customer.Role, customer.Surname, customer.Name))

		insertedCustomer, err := storage.Insert(context.Background(), customer)
		assert.NoError(t, err)

		// Check that the inserted customer matches what was expected
		expectedCustomer := models.Customer{
			Id:       id,
			Login:    customer.Login,
			Password: customer.Password,
			Role:     customer.Role,
			Surname:  customer.Surname,
			Name:     customer.Name,
		}
		assert.Equal(t, expectedCustomer, insertedCustomer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("InsertError", func(t *testing.T) {
		customer := models.Customer{
			Login:    "aisachenkov",
			Password: "Bebrolov",
			Role:     "admin",
			Surname:  "Isachenkov",
			Name:     "Artem",
		}

		mock.ExpectQuery("INSERT INTO "+os.Getenv("PSQL_TABLE_NAME")+
			" \\(login, password, role, surname, name\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id").
			WithArgs(customer.Login, customer.Password, customer.Role, customer.Surname, customer.Name).
			WillReturnError(fmt.Errorf("insert error"))

		insertedCustomer, err := storage.Insert(context.Background(), customer)
		assert.Error(t, err)
		assert.Nil(t, insertedCustomer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

// TestUpdate
/*
func TestUpdate(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("failed to open mock database: %s", err)
    }
    defer db.Close()

    os.Setenv("PSQL_TABLE_NAME", "customers") // Убедитесь, что имя таблицы установлено

    storage := MustNewPsqlCustomersStorage(db)

    var generatedID = uuid.New()

    t.Run("Success", func(t *testing.T) {
        customer := models.Customer{
            Id:       generatedID,
            Login:    "aisachenkov",
            Password: "Bebrolov",
            Role:     "admin",
            Surname:  "Isachenkov",
            Name:     "Artem",
        }

        // Установка ожидания для операции вставки
        mock.ExpectQuery("INSERT INTO "+os.Getenv("PSQL_TABLE_NAME")+
            " \\(login, password, role, surname, name\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id").
            WithArgs(customer.Login, customer.Password, customer.Role, customer.Surname, customer.Name).
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(customer.Id))

        // Вызов метода вставки
        _, err := storage.Insert(context.Background(), customer)
        assert.NoError(t, err)

        // Установка ожидания для получения клиента по ID перед обновлением
        mock.ExpectQuery("SELECT \\* FROM customers WHERE id = \\$1").
            WithArgs(customer.Id).
            WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "role", "surname", "name"}).
                AddRow(customer.Id, customer.Login, customer.Password, customer.Role, customer.Surname, customer.Name))

        // Обновление данных клиента
        updatedCustomer := models.Customer{
            Id:       customer.Id,
            Login:    "updated_login",
            Password: "updated_password",
            Role:     "admin",
            Surname:  "UpdatedSurname",
            Name:     "UpdatedName",
        }

        // Установка ожидания для операции обновления
        mock.ExpectExec("UPDATE "+os.Getenv("PSQL_TABLE_NAME")+" SET login = \\$1, password = \\$2, role = \\$3, surname = \\$4, name = \\$5 WHERE id = \\$6").
            WithArgs(updatedCustomer.Login, updatedCustomer.Password, updatedCustomer.Role, updatedCustomer.Surname, updatedCustomer.Name, updatedCustomer.Id).
            WillReturnResult(sqlmock.NewResult(1, 1))

        // Вызов метода обновления
        result, err := storage.Update(context.Background(), updatedCustomer)
        assert.NoError(t, err)

        // Проверка, что обновленный клиент соответствует ожиданиям
        assert.Equal(t, updatedCustomer, result)

        // Проверка выполнения всех ожиданий
        if err := mock.ExpectationsWereMet(); err != nil {
            t.Errorf("there were unfulfilled expectations: %s", err)
        }
    })

    t.Run("UpdateError", func(t *testing.T) {
        // Установка ожидания для получения клиента по ID
        mock.ExpectQuery("SELECT \\* FROM customers WHERE id = \\$1").
            WithArgs(generatedID).
            WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "role", "surname", "name"}).
                AddRow(generatedID, "aisachenkov", "Bebrolov", "admin", "Isachenkov", "Artem"))

        // Установка ожидания для операции обновления с ошибкой
        mock.ExpectExec("UPDATE "+os.Getenv("PSQL_TABLE_NAME")+" SET login = \\$1, password = \\$2, role = \\$3, surname = \\$4, name = \\$5 WHERE id = \\$6").
            WithArgs("NYC", "LAX", "admin", "UpdatedSurname", "UpdatedName", generatedID).
            WillReturnError(fmt.Errorf("update error"))

        _, err := storage.Update(context.Background(), models.Customer{
            Id:       generatedID,
            Login:    "updated_login",
            Password: "updated_password",
            Role:     "admin",
            Surname:  "UpdatedSurname",
            Name:     "UpdatedName",
        })
        assert.Error(t, err)

        // Проверка выполнения всех ожиданий
        if err := mock.ExpectationsWereMet(); err != nil {
            t.Errorf("there were unfulfilled expectations: %s", err)
        }
    })
}
*/

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlCustomersStorage(db)

	t.Run("Success", func(t *testing.T) {
		customerID := uuid.New()

		// Set up expectation for retrieving the customer
		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1").
			WithArgs(customerID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "role", "surname", "name"}).
				AddRow(customerID, "aisachenkov", "Bebrolov", "admin", "Isachenkov", "Artem"))

		// Set up expectation for deleting the customer
		mock.ExpectExec("DELETE FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1").
			WithArgs(customerID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		deletedCustomer, err := storage.Delete(context.Background(), customerID)
		assert.NoError(t, err)
		assert.NotNil(t, deletedCustomer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DeleteError", func(t *testing.T) {
		customerID := uuid.New()

		// Set up expectation for retrieving the customer
		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1").
			WithArgs(customerID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "password", "role", "surname", "name"}).
				AddRow(customerID, "aisachenkov", "Bebrolov", "admin", "Isachenkov", "Artem"))

		// Set up expectation for delete operation to fail
		mock.ExpectExec("DELETE FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1").
			WithArgs(customerID).
			WillReturnError(fmt.Errorf("delete error"))

		deletedCustomer, err := storage.Delete(context.Background(), customerID)
		assert.Error(t, err)
		assert.Nil(t, deletedCustomer)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
