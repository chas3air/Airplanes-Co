package psql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chas3air/Airplanes-Co/Core/DAL_flights/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlFlightsStorage(db)

	t.Run("Success", func(t *testing.T) {
		flightTime1, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
		flightTime2, _ := time.Parse(time.RFC3339, "2023-10-26T15:00:00Z")

		id1 := uuid.New()
		id2 := uuid.New()

		expectedFlights := []models.Flight{
			{Id: id1, FromWhere: "NYC", Destination: "LAX", FlightTime: flightTime1, FlightDuration: 6, FlightSeatsCosts: []int{100, 150, 200, 250}},
			{Id: id2, FromWhere: "LAX", Destination: "NYC", FlightTime: flightTime2, FlightDuration: 6, FlightSeatsCosts: []int{120, 170, 220, 270}},
		}

		rows := sqlmock.NewRows([]string{"id", "fromWhere", "destination", "flightTime", "flightDuration", "flightSeatsCost"}).
			AddRow(id1, expectedFlights[0].FromWhere, expectedFlights[0].Destination, expectedFlights[0].FlightTime, expectedFlights[0].FlightDuration, pq.Array(expectedFlights[0].FlightSeatsCosts)).
			AddRow(id2, expectedFlights[1].FromWhere, expectedFlights[1].Destination, expectedFlights[1].FlightTime, expectedFlights[1].FlightDuration, pq.Array(expectedFlights[1].FlightSeatsCosts))

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + ";").
			WillReturnRows(rows)

		flights, err := storage.GetAll(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedFlights, flights)
	})

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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlFlightsStorage(db)

	t.Run("Success", func(t *testing.T) {
		flightTime, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
		expectedId := uuid.New() // Генерируем уникальный UUID
		expectedFlight := models.Flight{
			Id:               expectedId,
			FromWhere:        "NYC",
			Destination:      "LAX",
			FlightTime:       flightTime,
			FlightDuration:   6,
			FlightSeatsCosts: []int{100, 150, 200, 250},
		}

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1;").
			WithArgs(expectedFlight.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fromWhere", "destination", "flightTime", "flightDuration", "flightSeatsCost"}).
				AddRow(expectedFlight.Id, expectedFlight.FromWhere, expectedFlight.Destination, expectedFlight.FlightTime, expectedFlight.FlightDuration, pq.Array(expectedFlight.FlightSeatsCosts)))

		flight, err := storage.GetById(context.Background(), expectedFlight.Id)
		assert.NoError(t, err)
		assert.Equal(t, expectedFlight, flight)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		// Генерируем новый UUID для этого теста
		notFoundId := uuid.New() // Новый UUID для теста "NotFound"
		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1;").
			WithArgs(notFoundId). // Используем notFoundId
			WillReturnError(sql.ErrNoRows)

		flight, err := storage.GetById(context.Background(), notFoundId) // Передаем notFoundId
		assert.Error(t, err)
		assert.Nil(t, flight)

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

	s := MustNewPsqlFlightsStorage(db)

	flightTime, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
	newFlight := models.Flight{
		FromWhere:        "NYC",
		Destination:      "LAX",
		FlightTime:       flightTime,
		FlightDuration:   6,
		FlightSeatsCosts: []int{100, 150, 200, 250},
	}

	mock.ExpectQuery("INSERT INTO "+os.Getenv("PSQL_TABLE_NAME")+" \\(fromWhere, destination, flightTime, flightDuration, flightSeatsCost\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\) RETURNING id").
		WithArgs(newFlight.FromWhere, newFlight.Destination, newFlight.FlightTime, newFlight.FlightDuration, pq.Array(newFlight.FlightSeatsCosts[:])).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uuid.New()))

	flight, err := s.Insert(context.Background(), newFlight)
	assert.NoError(t, err)
	assert.NotNil(t, flight)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlFlightsStorage(db)

	var generatedID = uuid.New()
	t.Run("Success", func(t *testing.T) {
		flightTime, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
		updatedFlight := models.Flight{
			Id:               generatedID,
			FromWhere:        "NYC",
			Destination:      "LAX",
			FlightTime:       flightTime,
			FlightDuration:   6,
			FlightSeatsCosts: []int{100, 150, 200, 250},
		}

		mock.ExpectExec("UPDATE "+os.Getenv("PSQL_TABLE_NAME")+" SET fromWhere = \\$1, destination = \\$2, flightTime = \\$3, flightDuration = \\$4, flightSeatsCost = \\$5 WHERE id = \\$6").
			WithArgs(updatedFlight.FromWhere, updatedFlight.Destination, updatedFlight.FlightTime, updatedFlight.FlightDuration, pq.Array(updatedFlight.FlightSeatsCosts), updatedFlight.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		flight, err := storage.Update(context.Background(), updatedFlight)
		assert.NoError(t, err)
		assert.NotNil(t, flight)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("UpdateError", func(t *testing.T) {
		mock.ExpectExec("UPDATE "+os.Getenv("PSQL_TABLE_NAME")+" SET fromWhere = \\$1, destination = \\$2, flightTime = \\$3, flightDuration = \\$4, flightSeatsCost = \\$5 WHERE id = \\$6").
			WithArgs("NYC", "LAX", time.Now(), 6, pq.Array([]int{100, 150, 200, 250}), generatedID). // Используем срез
			WillReturnError(fmt.Errorf("update error"))

		_, err := storage.Update(context.Background(), models.Flight{
			Id:               generatedID,
			FromWhere:        "NYC",
			Destination:      "LAX",
			FlightDuration:   6,
			FlightTime:       time.Now(),
			FlightSeatsCosts: []int{100, 150, 200, 250},
		})
		assert.Error(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %s", err)
	}
	defer db.Close()

	storage := MustNewPsqlFlightsStorage(db)

	t.Run("Success", func(t *testing.T) {
		flightTime, _ := time.Parse(time.RFC3339, "2023-10-25T15:00:00Z")
		flightToDelete := models.Flight{
			Id:               uuid.New(),
			FromWhere:        "NYC",
			Destination:      "LAX",
			FlightTime:       flightTime,
			FlightDuration:   6,
			FlightSeatsCosts: []int{100, 150, 200, 250},
		}

		mock.ExpectQuery("SELECT \\* FROM " + os.Getenv("PSQL_TABLE_NAME") + " WHERE id = \\$1;").
			WithArgs(flightToDelete.Id).
			WillReturnRows(sqlmock.NewRows([]string{"id", "fromWhere", "destination", "flightTime", "flightDuration", "flightSeatsCost"}).
				AddRow(flightToDelete.Id, flightToDelete.FromWhere, flightToDelete.Destination, flightToDelete.FlightTime, flightToDelete.FlightDuration, pq.Array(flightToDelete.FlightSeatsCosts[:])))

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
