package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/chas3air/Airplanes-Co/DAL/internal/config"
	"github.com/chas3air/Airplanes-Co/DAL/internal/models"

	_ "github.com/lib/pq"
)

type PsqlTicketsStorage struct {
	PsqlStorage
}

func MustNewPsqlTicketsStorage(db *sql.DB) PsqlTicketsStorage {
	const op = "DAL.internal.storage.psqlRepository.newTicketStorage"
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ` + config.PSQL_TICKETS_TABLE_NAME + ` (
			id SERIAL PRIMARY KEY,
			flightFromWhere VARCHAR(100) NOT NULL,
			flightDestination VARCHAR(100) NOT NULL,
			flightTime TIMESTAMP NOT NULL,
			ownerId INT NOT NULL,
			ticketCost NUMERIC(10, 2) NOT NULL,
			classOfService VARCHAR(50) NOT NULL
		);
	`)

	if err != nil {
		panic("Ошибка при создании таблицы для пользователей: " + fmt.Errorf("%s: %s", op, err).Error())
	}

	return PsqlTicketsStorage{
		PsqlStorage: NewPsqlStorage(db),
	}
}

func (s PsqlTicketsStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.GetAll"

	rows, err := s.DB.QueryContext(ctx, `SELECT * FROM `+config.PSQL_TICKETS_TABLE_NAME+`;`)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()
	var ticket models.Ticket
	tickets := make([]models.Ticket, 0, 10)

	for rows.Next() {
		err := rows.Scan(&ticket.Id, &ticket.FlightFromWhere,
			&ticket.FlightDestination, &ticket.FlightTime, &ticket.OwnerId,
			&ticket.TicketCost, &ticket.ClassOfService)

		if err != nil {
			log.Println("Строка несчитана, ошибка: ", err.Error())
			continue
		}
		tickets = append(tickets, ticket)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return tickets, nil
}

func (s PsqlTicketsStorage) GetById(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.GetById"

	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+config.PSQL_TICKETS_TABLE_NAME+`
		WHERE id=$1;
	`, id)

	var ticket models.Ticket

	err := row.Scan(&ticket.Id, &ticket.FlightFromWhere, &ticket.FlightDestination,
		&ticket.FlightTime, &ticket.OwnerId, &ticket.TicketCost, &ticket.ClassOfService)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: запись с id=%d не найдена", op, id)
		}
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	return ticket, nil
}

func (s PsqlTicketsStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Insert"

	ticket := innerObj.(models.Ticket)
	var id int = -1

	err := s.DB.QueryRowContext(ctx, `
		INSERT INTO `+config.PSQL_TICKETS_TABLE_NAME+`
		(flightFromWhere, flightDestination, flightTime, ownerId, ticketCost, classOfServices)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
		WHERE id = $7
	`, ticket.FlightFromWhere, ticket.FlightDestination, ticket.FlightTime, ticket.OwnerId, ticket.TicketCost, ticket.ClassOfService, ticket.Id)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return s.GetById(ctx, id)
}

func (s PsqlTicketsStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Update"

	ticket := innerObj.(models.Ticket)

	_, err := s.DB.ExecContext(ctx, `
		UPDATE `+config.PSQL_TICKETS_TABLE_NAME+`
		SET flightFromWhere = $1, flightDestination = $2, flightTime = $3, ownerId = $4, ticketCost = $5, classOfServices = $6
		WHERE id = $7;
	`, ticket.FlightFromWhere, ticket.FlightDestination, ticket.FlightTime, ticket.OwnerId, ticket.TicketCost, ticket.ClassOfService, ticket.Id)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return s.GetById(ctx, ticket.Id)
}

func (s PsqlTicketsStorage) Delete(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Delete"

	ticket, err := s.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+config.PSQL_TICKETS_TABLE_NAME+`
		WHERE id = $1;
	`, id)

	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return ticket, nil
}
