package psql

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/models"
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
