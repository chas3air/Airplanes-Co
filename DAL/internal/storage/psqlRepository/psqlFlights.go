package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/chas3air/Airplanes-Co/DAL/internal/config"
	"github.com/chas3air/Airplanes-Co/DAL/internal/models"
)

type PsqlFlightsStorage struct {
	PsqlStorage
}

func MustNewPsqlFlightsStorage(db *sql.DB) PsqlFlightsStorage {
	const op = "DAL.internal.storage.psqlRepository.newFlightsStorage"
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + config.PSQL_FLIGHTS_TABLE_NAME + ` (
			id SERIAL PRIMARY KEY,
			from_where VARCHAR(50) NOT NULL,
			destination VARCHAR(50) NOT NULL,
			flight_time TIMESTAMP(50) NOT NULL,
			flight_duration INT NOT NULL
		);
	`)

	if err != nil {
		panic("Ошибка при создании таблицы для рейсов: " + fmt.Errorf("%s: %w", op, err).Error())
	}

	return PsqlFlightsStorage{
		PsqlStorage: NewPsqlStorage(db),
	}
}

func (s PsqlFlightsStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.GetAll"
	rows, err := s.DB.QueryContext(ctx, `SELECT * FROM `+config.PSQL_FLIGHTS_TABLE_NAME+`;`)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	defer rows.Close()
	flights := make([]models.Flight, 0, 10)
	var flight models.Flight

	for rows.Next() {
		err := rows.Scan(&flight.Id, &flight.FromWhere, &flight.Destination, &flight.FlightTime, &flight.FlightDuration)
		if err != nil {
			log.Println("Строка несчитана, ошибка", err.Error())
			continue
		}
		flights = append(flights, flight)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return flights, nil
}

func (s PsqlFlightsStorage) GetById(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.GetById"
	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+config.PSQL_FLIGHTS_TABLE_NAME+`
		WHERE id =$1;
		`, id)

	var flight models.Flight
	err := row.Scan(&flight.Id, &flight.FromWhere, &flight.Destination, &flight.FlightTime, &flight.FlightDuration)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: запись с id=%d не найдена", op, id)
		}
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return flight, nil
}

func (s PsqlFlightsStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Insert"
	flight := innerObj.(models.Flight)
	var id int

	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO `+config.PSQL_FLIGHTS_TABLE_NAME+`
		(fromWhere, destination, flightTime, flightDuration)
		VALUES ($1, $2, $3, $4)  RETURNING id
		`, flight.FromWhere, flight.Destination, flight.FlightTime, flight.FlightDuration).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return s.GetById(ctx, id)
}

func (s PsqlFlightsStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Update"
	flight := innerObj.(models.Flight)
	_, err := s.DB.ExecContext(ctx, `
		UPDATE `+config.PSQL_FLIGHTS_TABLE_NAME+`
		SET fromWhere=$1, destination=$2, flightTime=$3, flightDuration=$4
		WHERE id=$5;
	`, flight.FromWhere, flight.Destination, flight.FlightTime, flight.FlightDuration, flight.Id)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return s.GetById(ctx, flight.Id)
}

func (s PsqlFlightsStorage) Delete(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlFlights.Delete"

	flight, err := s.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+config.PSQL_FLIGHTS_TABLE_NAME+`
		WHERE id = $1;
	`, id)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return flight, nil
}
