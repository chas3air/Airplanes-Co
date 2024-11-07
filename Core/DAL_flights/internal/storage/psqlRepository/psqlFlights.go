package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chas3air/Airplanes-Co/Core/DAL_flights/internal/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PsqlFlightsStorage struct {
	PsqlStorage
}

// MustNewPsqlFlightsStorage initializes a new PsqlFlightsStorage instance.
func MustNewPsqlFlightsStorage(db *sql.DB) PsqlFlightsStorage {
	return PsqlFlightsStorage{
		PsqlStorage: NewPsqlStorage(db),
	}
}

// GetAll retrieves all flights from the database.
func (s PsqlFlightsStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.GetAll"

	rows, err := s.DB.QueryContext(ctx, `SELECT * FROM `+os.Getenv("PSQL_TABLE_NAME")+`;`)
	if err != nil {
		log.Println("Error querying flights:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	flights := make([]models.Flight, 0)
	for rows.Next() {
		var flight models.Flight
		var seatCosts pq.Int64Array // Используем pq.Int64Array

		if err := rows.Scan(&flight.Id, &flight.FromWhere, &flight.Destination, &flight.FlightTime, &flight.FlightDuration, &seatCosts); err != nil {
			log.Println("Error scanning row:", err.Error())
			continue
		}

		// Проверка длины массива seatCosts
		if len(seatCosts) != 4 {
			log.Println("Invalid number of seat costs, expected 4 but got:", len(seatCosts))
			continue
		}

		// Инициализация массива
		flight.FlightSeatsCosts = make([]int, len(seatCosts))
		for i, cost := range seatCosts {
			flight.FlightSeatsCosts[i] = int(cost)
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
func (s PsqlFlightsStorage) GetById(ctx context.Context, id any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.GetById"
	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+os.Getenv("PSQL_TABLE_NAME")+`
		WHERE id = $1;
	`, id)

	var flight models.Flight
	var seatCosts pq.Int64Array // Используем pq.Int64Array
	if err := row.Scan(&flight.Id, &flight.FromWhere, &flight.Destination, &flight.FlightTime, &flight.FlightDuration, &seatCosts); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No flight found with ID=%v\n", id)
			return nil, fmt.Errorf("%s: no flight found with ID=%v", op, id)
		}
		log.Println("Error scanning flight:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	// Проверка длины массива seatCosts
	if len(seatCosts) != 4 {
		log.Println("Invalid number of seat costs, expected 4 but got:", len(seatCosts))
		return nil, fmt.Errorf("%s: invalid seat costs", op)
	}

	// Инициализация массива
	flight.FlightSeatsCosts = make([]int, len(seatCosts))
	for i, cost := range seatCosts {
		flight.FlightSeatsCosts[i] = int(cost)
	}

	log.Printf("Retrieved flight with ID=%s\n", flight.Id.String())
	return flight, nil
}

// Insert adds a new flight to the database.
func (s PsqlFlightsStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Insert"
	flight := innerObj.(models.Flight)

	var id uuid.UUID
	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO `+os.Getenv("PSQL_TABLE_NAME")+`
		(fromWhere, destination, flightTime, flightDuration, flightSeatsCost)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, flight.FromWhere, flight.Destination, flight.FlightTime, flight.FlightDuration, pq.Array(flight.FlightSeatsCosts[:])).Scan(&id)
	if err != nil {
		log.Println("Error inserting flight:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	flight.Id = id

	log.Printf("Inserted flight with ID=%s\n", id.String())
	return flight, nil
}

func (s PsqlFlightsStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Update"
	flight := innerObj.(models.Flight)

	_, err := s.DB.ExecContext(ctx, `
		UPDATE `+os.Getenv("PSQL_TABLE_NAME")+`
		SET fromWhere = $1, destination = $2, flightTime = $3, flightDuration = $4, flightSeatsCost = $5
		WHERE id = $6;
	`, flight.FromWhere, flight.Destination, flight.FlightTime, flight.FlightDuration, pq.Array(flight.FlightSeatsCosts), flight.Id)

	if err != nil {
		log.Println("Error updating flight:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Updated flight with ID=%s\n", flight.Id.String())
	return flight, nil
}

// Delete removes a flight from the database by its ID.
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
