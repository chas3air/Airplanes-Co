package psql

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/chas3air/Airplanes-Co/Core/DAL_airplanes/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Create a PsqlAirplanesStorage instance with the mocked DB
	storage := PsqlAirplanesStorage{PsqlStorage: PsqlStorage{DB: db}}

	// Mocking the query
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(uuid.New(), "Boeing 737").
		AddRow(uuid.New(), "Airbus A320")

	mock.ExpectQuery("SELECT id, name FROM").WillReturnRows(rows)

	// Execute the method
	airplanes, err := storage.GetAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, airplanes, 2) // Ensure this matches the expected count
}

func TestInsert(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := PsqlAirplanesStorage{PsqlStorage: PsqlStorage{DB: db}}

	airplane := models.Airplane{Id: uuid.New(), Name: "Boeing 737"}

	// Mocking the insert
	mock.ExpectExec("INSERT INTO").
		WithArgs(airplane.Id, airplane.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Execute the method
	result, err := storage.Insert(context.Background(), airplane)
	assert.NoError(t, err)
	assert.Equal(t, airplane, result)
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := PsqlAirplanesStorage{PsqlStorage: PsqlStorage{DB: db}}

	airplane := models.Airplane{Id: uuid.New(), Name: "Boeing 737"}

	// Mocking the update
	mock.ExpectExec("UPDATE").
		WithArgs(airplane.Name, airplane.Id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the method
	result, err := storage.Update(context.Background(), airplane)
	assert.NoError(t, err)
	assert.Equal(t, airplane, result)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := PsqlAirplanesStorage{PsqlStorage: PsqlStorage{DB: db}}

	airplaneID := uuid.New()
	airplane := models.Airplane{Id: airplaneID, Name: "Boeing 737"}

	// Mocking the select before delete
	mock.ExpectQuery("SELECT \\* FROM").
		WithArgs(airplaneID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(airplane.Id, airplane.Name))

	// Mocking the delete
	mock.ExpectExec("DELETE FROM").
		WithArgs(airplaneID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Execute the method
	result, err := storage.Delete(context.Background(), airplaneID)
	assert.NoError(t, err)
	assert.Equal(t, airplane, result)
}
