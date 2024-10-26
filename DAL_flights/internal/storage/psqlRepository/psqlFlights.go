package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chas3air/Airplanes-Co/DAL_flights/internal/models"
	"github.com/google/uuid"
)

type PsqlFlightsStorage struct {
	PsqlStorage
}

// MustNewPsqlFlightsStorage initializes a new PsqlFlightsStorage instance.
// It checks the database connection and panics if the flights table is unavailable.
func MustNewPsqlFlightsStorage(db *sql.DB) PsqlFlightsStorage {
	const op = "DAL.internal.storage.psqlRepository.newFlightsStorage"
	err := db.Ping()
	if err != nil {
		log.Println("Flights table is unavailable: " + err.Error())
		log.Panic(fmt.Errorf("%s: %w", op, err))
	}

	return PsqlFlightsStorage{
		PsqlStorage: NewPsqlStorage(db),
	}
}

// GetAll retrieves all flights from the database.
// It returns a slice of Flight models and an error, if any occurs.
func (s PsqlFlightsStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.GetAll"
	rows, err := s.DB.QueryContext(ctx, `SELECT * FROM `+os.Getenv("PSQL_TABLE_NAME")+`;`)
	if err != nil {
		log.Println("Error querying flights:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	flights := make([]models.Flight, 0, 10)
	var flight models.Flight

	for rows.Next() {
		err := rows.Scan(&flight.Id, &flight.FromWhere, &flight.Destination, &flight.FlightTime, &flight.FlightDuration)
		if err != nil {
			log.Println("Error scanning row:", err.Error())
			continue
		}
		flights = append(flights, flight)
	}

	if err := rows.Err(); err != nil {
		log.Println("Row error:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrieved %d flights\n", len(flights))
	return flights, nil
}

// GetById retrieves a flight by its ID.
// It returns the Flight model and an error if no flight is found or if an error occurs.
func (s PsqlFlightsStorage) GetById(ctx context.Context, id any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.GetById"
	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+os.Getenv("PSQL_TABLE_NAME")+`
		WHERE id = $1;
	`, id)

	var flight models.Flight
	err := row.Scan(&flight.Id, &flight.FromWhere, &flight.Destination, &flight.FlightTime, &flight.FlightDuration)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No flight found with ID=%v\n", id)
			return nil, fmt.Errorf("%s: no flight found with ID=%v", op, id)
		}
		log.Println("Error scanning flight:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrieved flight with ID=%d\n", flight.Id)
	return flight, nil
}

// Insert adds a new flight to the database.
// It returns the inserted Flight model and an error if the insertion fails.
func (s PsqlFlightsStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Insert"
	flight := innerObj.(models.Flight)
	var id uuid.UUID

	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO `+os.Getenv("PSQL_TABLE_NAME")+`
		(fromWhere, destination, flightTime, flightDuration)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, flight.FromWhere, flight.Destination, flight.FlightTime, flight.FlightDuration).Scan(&id)
	if err != nil {
		log.Println("Error inserting flight:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	flight.Id = id

	log.Printf("Inserted flight with ID=%d\n", id)
	return flight, nil
}

// Update modifies an existing flight in the database.
// It returns the updated Flight model and an error if the update fails.
func (s PsqlFlightsStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Update"
	flight := innerObj.(models.Flight)

	_, err := s.DB.ExecContext(ctx, `
		UPDATE `+os.Getenv("PSQL_TABLE_NAME")+`
		SET fromWhere = $1, destination = $2, flightTime = $3, flightDuration = $4
		WHERE id = $5;
	`, flight.FromWhere, flight.Destination, flight.FlightTime, flight.FlightDuration, flight.Id)

	if err != nil {
		log.Println("Error updating flight:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Updated flight with ID=%d\n", flight.Id)
	return flight, nil
}

// Delete removes a flight from the database by its ID.
// It returns the deleted Flight model and an error if the deletion fails.
func (s PsqlFlightsStorage) Delete(ctx context.Context, id any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Delete"

	flight, err := s.GetById(ctx, id)
	if err != nil {
		log.Println("Error getting flight by ID:", err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+os.Getenv("PSQL_TABLE_NAME")+`
		WHERE id = $1;
	`, id)

	if err != nil {
		log.Println("Error deleting flight:", err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Deleted flight with ID=%d\n", id)
	return flight, nil
}
