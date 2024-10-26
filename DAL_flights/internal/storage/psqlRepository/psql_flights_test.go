package psql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chas3air/Airplanes-Co/DAL_flights/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	// Mock database setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlFlightsStorage(db)

	// Test for successful retrieval of all flights
	t.Run("Success", func(t *testing.T) {
		flightTime1, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
		flightTime2, _ := time.Parse(time.RFC3339, "2023-10-26T15:00:00Z")

		expectedFlights := []models.Flight{
			{Id: uuid.New(), FromWhere: "NYC", Destination: "LAX", FlightTime: flightTime1, FlightDuration: 6},
			{Id: uuid.New(), FromWhere: "LAX", Destination: "NYC", FlightTime: flightTime2, FlightDuration: 6},
		}

		rows := sqlmock.NewRows([]string{"id", "fromWhere", "destination", "flightTime", "flightDuration"}).
			AddRow(expectedFlights[0].Id, expectedFlights[0].FromWhere, expectedFlights[0].Destination, expectedFlights[0].FlightTime, expectedFlights[0].FlightDuration).
			AddRow(expectedFlights[1].Id, expectedFlights[1].FromWhere, expectedFlights[1].Destination, expectedFlights[1].FlightTime, expectedFlights[1].FlightDuration)

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + ";").
			WillReturnRows(rows)

		flights, err := storage.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedFlights, flights)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// Test for error in querying
	t.Run("QueryError", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + ";").
			WillReturnError(fmt.Errorf("query error"))

		flights, err := storage.GetAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, flights)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestGetById(t *testing.T) {
	// Mock database setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlFlightsStorage(db)

	// Test for successful retrieval by ID
	t.Run("Success", func(t *testing.T) {
		flightTime, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
		expectedFlight := models.Flight{
			Id:             uuid.New(),
			FromWhere:      "NYC",
			Destination:    "LAX",
			FlightTime:     flightTime,
			FlightDuration: 6,
		}

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1;").
			WithArgs(expectedFlight.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "from_where", "destination", "flight_time", "flight_duration"}).
				AddRow(expectedFlight.Id, expectedFlight.FromWhere, expectedFlight.Destination, expectedFlight.FlightTime, expectedFlight.FlightDuration))

		flight, err := storage.GetById(context.Background(), expectedFlight.Id)
		assert.NoError(t, err)
		assert.Equal(t, expectedFlight, flight)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// Test for flight not found
	t.Run("NotFound", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1;").
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)

		flight, err := storage.GetById(context.Background(), 2)
		assert.Error(t, err)
		assert.Nil(t, flight)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestInsert(t *testing.T) {
	// Настройка мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	s := MustNewPsqlFlightsStorage(db)

	// Создание нового рейса для вставки
	flightTime, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
	newFlight := models.Flight{
		FromWhere:      "NYC",
		Destination:    "LAX",
		FlightTime:     flightTime,
		FlightDuration: 6,
	}

	// Ожидание выполнения вставки
	mock.ExpectQuery("INSERT INTO "+os.Getenv("PSQL_TABLE_NAME")+" \\(fromWhere, destination, flightTime, flightDuration\\) VALUES \\(\\$1, \\$2, \\$3, \\$4\\) RETURNING id").
		WithArgs(newFlight.FromWhere, newFlight.Destination, newFlight.FlightTime, newFlight.FlightDuration).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New())) // Возвращаем новый UUID

	// Вызов метода вставки
	flight, err := s.Insert(context.Background(), newFlight)
	assert.NoError(t, err)
	assert.NotNil(t, flight)

	// Проверка выполнения всех ожиданий
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	// Mock database setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlFlightsStorage(db)

	var generatedID = uuid.New()
	// Test for successful update
	t.Run("Success", func(t *testing.T) {
		flightTime, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
		updatedFlight := models.Flight{
			Id:             generatedID,
			FromWhere:      "NYC",
			Destination:    "LAX",
			FlightTime:     flightTime,
			FlightDuration: 6,
		}

		// Ожидание выполнения запроса на обновление
		mock.ExpectExec("UPDATE "+os.Getenv("PSQL_TABLE_NAME")+" SET fromWhere = \\$1, destination = \\$2, flightTime = \\$3, flightDuration = \\$4 WHERE id = \\$5").
			WithArgs(updatedFlight.FromWhere, updatedFlight.Destination, updatedFlight.FlightTime, updatedFlight.FlightDuration, updatedFlight.Id).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Возвращаем результат обновления

		flight, err := storage.Update(context.Background(), updatedFlight)
		assert.NoError(t, err)
		assert.NotNil(t, flight)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// Test for update error
	t.Run("UpdateError", func(t *testing.T) {
		mock.ExpectExec("UPDATE "+os.Getenv("PSQL_TABLE_NAME")+" SET fromWhere = \\$1, destination = \\$2, flightTime = \\$3, flightDuration = \\$4 WHERE id = \\$5").
			WithArgs("NYC", "LAX", time.Now(), 6, generatedID).
			WillReturnError(fmt.Errorf("update error"))

		_, err := storage.Update(context.Background(), models.Flight{Id: generatedID, FromWhere: "NYC", Destination: "LAX", FlightDuration: 6, FlightTime: time.Now()})
		assert.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
func TestDelete(t *testing.T) {
	// Mock database setup
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlFlightsStorage(db)

	// Test for successful deletion
	t.Run("Success", func(t *testing.T) {
		flightTime, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
		flightToDelete := models.Flight{
			Id:             uuid.New(),
			FromWhere:      "NYC",
			Destination:    "LAX",
			FlightTime:     flightTime,
			FlightDuration: 6,
		}

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1;").
			WithArgs(flightToDelete.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "from_where", "destination", "flight_time", "flight_duration"}).
				AddRow(flightToDelete.Id, flightToDelete.FromWhere, flightToDelete.Destination, flightToDelete.FlightTime, flightToDelete.FlightDuration))

		mock.ExpectExec("DELETE FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1;").
			WithArgs(flightToDelete.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		deletedFlight, err := storage.Delete(context.Background(), flightToDelete.Id)
		assert.NoError(t, err)
		assert.Equal(t, flightToDelete, deletedFlight)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	// Test for delete error
	t.Run("DeleteError", func(t *testing.T) {
		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1;").
			WithArgs(2).
			WillReturnError(fmt.Errorf("no flight found"))

		_, err := storage.Delete(context.Background(), 2)
		assert.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
