package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chas3air/Airplanes-Co/Core/DAL_flights/internal/models"
	"github.com/google/uuid"
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
// GetAll retrieves all flights from the database.
func (s PsqlFlightsStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.GetAll"

	rows, err := s.DB.QueryContext(ctx, `
		SELECT f.id, c1.name AS fromWhere, c2.name AS destination, f.flightTime, f.flightDuration, f.airplaneId
		FROM Flights f
		JOIN City c1 ON f.fromWhere = c1.id
		JOIN City c2 ON f.destination = c2.id;
	`)
	if err != nil {
		log.Printf("%s: error querying flights: %v\n", op, err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	var flights []models.Flight
	for rows.Next() {
		var flight models.Flight

		if err := rows.Scan(&flight.Id, &flight.FromWhere, &flight.Destination, &flight.FlightTime, &flight.FlightDuration, &flight.Airplane); err != nil {
			log.Printf("%s: error scanning row: %v\n", op, err)
			continue
		}

		// Получаем стоимость мест для текущего рейса
		seatCosts, err := s.getSeatCostsByFlightId(ctx, flight.Id)
		if err == nil {
			flight.FlightSeatsCosts = seatCosts
		}

		flights = append(flights, flight)
	}

	if err := rows.Err(); err != nil {
		log.Printf("%s: row error: %v\n", op, err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrieved %d flights\n", len(flights))
	return flights, nil
}

// getSeatCostsByFlightId retrieves seat costs for a given flight ID.
// It ensures that the returned slice always contains 4 values: Economy, Comfort, Business, First Class.
func (s PsqlFlightsStorage) getSeatCostsByFlightId(ctx context.Context, flightId uuid.UUID) ([]int, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.getSeatCostsByFlightId"

	costs := make([]int, 4)

	rows, err := s.DB.QueryContext(ctx, `
		SELECT className, cost FROM CostTable WHERE flightId = $1;
	`, flightId)
	if err != nil {
		log.Printf("%s: error querying seat costs: %v\n", op, err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var className string
		var cost int
		if err := rows.Scan(&className, &cost); err != nil {
			log.Printf("%s: error scanning seat cost: %v\n", op, err)
			continue
		}

		switch className {
		case "Economy":
			costs[0] = cost
		case "Comfort":
			costs[1] = cost
		case "Business":
			costs[2] = cost
		case "First Class":
			costs[3] = cost
		}
	}

	if err := rows.Err(); err != nil {
		log.Printf("%s: row error: %v\n", op, err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return costs, nil
}

// GetById retrieves a flight by its ID.
// GetById retrieves a flight by its ID, including seat costs.
func (s PsqlFlightsStorage) GetById(ctx context.Context, id any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.GetById"

	row := s.DB.QueryRowContext(ctx, `
		SELECT f.id, c1.name AS fromWhere, c2.name AS destination, f.flightTime, f.flightDuration, f.airplaneId
		FROM Flights f
		JOIN City c1 ON f.fromWhere = c1.id
		JOIN City c2 ON f.destination = c2.id
		WHERE f.id = $1;
	`, id)

	var flight models.Flight

	if err := row.Scan(&flight.Id, &flight.FromWhere, &flight.Destination, &flight.FlightTime, &flight.FlightDuration, &flight.Airplane); err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No flight found with ID=%v\n", id)
			return nil, fmt.Errorf("%s: no flight found with ID=%v", op, id)
		}
		log.Println("Error scanning flight:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	// Получаем стоимость мест для текущего рейса
	seatCosts, err := s.getSeatCostsByFlightId(ctx, flight.Id)
	if err == nil {
		flight.FlightSeatsCosts = seatCosts
	}

	log.Printf("Retrieved flight with ID=%s\n", flight.Id.String())
	return flight, nil
}

// Insert adds a new flight to the database, including seat costs.
func (s PsqlFlightsStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Insert"
	flight := innerObj.(models.Flight)

	row := s.DB.QueryRowContext(ctx, `
		SELECT c.id From City c
		WHERE name = $1;
	`, flight.FromWhere)
	var from int
	if err := row.Scan(&from); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No city found with name=", flight.FromWhere)
			return nil, fmt.Errorf("%s: no city found with name=%s", op, flight.FromWhere)
		}
	}

	row = s.DB.QueryRowContext(ctx, `
		SELECT c.id From City c
		WHERE name = $1;
	`, flight.Destination)
	var dest int
	if err := row.Scan(&dest); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No city found with name=", flight.FromWhere)
			return nil, fmt.Errorf("%s: no city found with name=%s", op, flight.FromWhere)
		}
	}

	// Вставляем рейс в таблицу Flights
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO Flights
		(id, fromWhere, destination, flightTime, flightDuration, airplaneId)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, flight.Id, from, dest, flight.FlightTime, flight.FlightDuration, flight.Airplane)
	if err != nil {
		log.Printf("%s: error inserting flight: %v\n", op, err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	// Вставляем стоимости мест в таблицу CostTable
	classNames := []string{"Economy", "Comfort", "Business", "First Class"}
	for i, cost := range flight.FlightSeatsCosts {
		_, err := s.DB.ExecContext(ctx, `
			INSERT INTO CostTable (flightId, className, cost)
			VALUES ($1, $2, $3)
		`, flight.Id, classNames[i], cost)
		if err != nil {
			log.Printf("%s: error inserting seat cost for class %s: %v\n", op, classNames[i], err)
			return nil, fmt.Errorf("%s: %v", op, err)
		}
	}

	log.Printf("Inserted flight with ID=%s and associated seat costs\n", flight.Id.String())
	return flight, nil
}

// Update modifies an existing flight in the database, including seat costs.
func (s PsqlFlightsStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Update"
	flight := innerObj.(models.Flight)

	row := s.DB.QueryRowContext(ctx, `
		SELECT c.id From City c
		WHERE name = $1;
	`, flight.FromWhere)
	var from int
	if err := row.Scan(&from); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No city found with name=", flight.FromWhere)
			return nil, fmt.Errorf("%s: no city found with name=%s", op, flight.FromWhere)
		}
	}

	row = s.DB.QueryRowContext(ctx, `
		SELECT c.id From City c
		WHERE name = $1;
	`, flight.Destination)
	var dest int
	if err := row.Scan(&dest); err != nil {
		if err == sql.ErrNoRows {
			log.Println("No city found with name=", flight.FromWhere)
			return nil, fmt.Errorf("%s: no city found with name=%s", op, flight.FromWhere)
		}
	}

	// Обновляем информацию о рейсе в таблице Flights
	_, err := s.DB.ExecContext(ctx, `
		UPDATE Flights
		SET fromWhere = $1, destination = $2, flightTime = $3, flightDuration = $4, airplaneId = $5
		WHERE id = $6;
	`, from, dest, flight.FlightTime, flight.FlightDuration, flight.Airplane, flight.Id)

	if err != nil {
		log.Printf("%s: error updating flight: %v\n", op, err)
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	// Обновляем стоимости мест в таблице CostTable
	classNames := []string{"Economy", "Comfort", "Business", "First Class"}
	for i, cost := range flight.FlightSeatsCosts {
		_, err := s.DB.ExecContext(ctx, `
			INSERT INTO CostTable (flightId, className, cost)
			VALUES ($1, $2, $3)
			ON CONFLICT (flightId, className) DO UPDATE SET cost = EXCLUDED.cost;
		`, flight.Id, classNames[i], cost)
		if err != nil {
			log.Printf("%s: error updating seat cost for class %s: %v\n", op, classNames[i], err)
			return nil, fmt.Errorf("%s: %v", op, err)
		}
	}

	log.Printf("Updated flight with ID=%s and associated seat costs\n", flight.Id.String())
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
		DELETE FROM CostTable
		WHERE flightId = $1;
	`, id)
	if err != nil {
		log.Println("Error deleting seat costs associated with flight:", err)
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

	log.Printf("Deleted flight with ID=%s\n", id)
	return flight, nil
}
