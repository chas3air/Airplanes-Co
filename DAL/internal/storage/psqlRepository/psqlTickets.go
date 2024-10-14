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

// TODO: все переделать запросы

func MustNewPsqlTicketsStorage(db *sql.DB) PsqlTicketsStorage {
	const op = "DAL.internal.storage.psqlRepository.newTicketStorage"
	err := db.Ping()
	if err != nil {
		log.Println("Table of tickets is unavailable: " + err.Error())
		panic("Ошибка при создании таблицы для пользователей: " + fmt.Errorf("%s: %s", op, err).Error())
	}

	return PsqlTicketsStorage{
		PsqlStorage: NewPsqlStorage(db),
	}
}

func (s PsqlTicketsStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.GetAll"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of tickets is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of tickets is unavailable, file: %s: %w", op, err)
	}

	rows, err := s.DB.QueryContext(ctx, `SELECT * FROM `+config.PSQL_TICKETS_TABLE_NAME+`;`)
	if err != nil {
		log.Println("Error querying customers:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	var ticket models.Ticket
	tickets := make([]models.Ticket, 0, 10)

	for rows.Next() {
		err := rows.Scan(&ticket.Id, &ticket.FlightId, &ticket.OwnerId,
			&ticket.TicketCost, &ticket.ClassOfService)
		if err != nil {
			log.Println("Error scanning row:", err.Error())
			continue
		}
		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		log.Println("Row error:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Println("Count of retrived tickets is:", len(tickets))
	return tickets, nil
}

func (s PsqlTicketsStorage) GetById(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.GetById"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of tickets is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of tickets is unavailable, file: %s: %w", op, err)
	}

	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+config.PSQL_TICKETS_TABLE_NAME+`
		WHERE id=$1;
	`, id)

	var ticket models.Ticket

	err = row.Scan(&ticket.Id, &ticket.FlightId, &ticket.OwnerId, &ticket.TicketCost, &ticket.ClassOfService)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No ticket found with id=%d\n", id)
			return nil, fmt.Errorf("%s: No ticket found with id=%d", op, id)
		}
		log.Println("Error scanning ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrived ticket`s id=%d\n", ticket.Id)
	return ticket, nil
}

func (s PsqlTicketsStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Insert"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of tickets is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of tickets is unavailable, file: %s: %w", op, err)
	}

	ticket := innerObj.(models.Ticket)
	var id int = -1

	err = s.DB.QueryRowContext(ctx, `
		INSERT INTO `+config.PSQL_TICKETS_TABLE_NAME+`
		(flightId, ownerId, ticketCost, classOfServices)
		VALUES ($1, $2, $3, $4) RETURNING id
		WHERE id = $5
	`, ticket.FlightId, ticket.OwnerId, ticket.TicketCost, ticket.ClassOfService, ticket.Id).Scan(&id)

	if err != nil {
		log.Println("Error inserting ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return s.GetById(ctx, id)
}

func (s PsqlTicketsStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Update"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of tickets is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of tickets is unavailable, file: %s: %w", op, err)
	}

	ticket := innerObj.(models.Ticket)

	_, err = s.DB.ExecContext(ctx, `
		UPDATE `+config.PSQL_TICKETS_TABLE_NAME+`
		SET flightId = $1, ownerId = $2, ticketCost = $3, classOfServices = $4
		WHERE id = $5;
	`, ticket.FlightId, ticket.OwnerId, ticket.TicketCost, ticket.ClassOfService, ticket.Id)

	if err != nil {
		log.Println("Error updating customer:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return s.GetById(ctx, ticket.Id)
}

func (s PsqlTicketsStorage) Delete(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Delete"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Table of tickets is unavailable: " + err.Error())
		return nil, fmt.Errorf("table of tickets is unavailable, file: %s: %w", op, err)
	}

	ticket, err := s.GetById(ctx, id)
	if err != nil {
		log.Println("Error getting ticket by id:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+config.PSQL_TICKETS_TABLE_NAME+`
		WHERE id = $1;
	`, id)

	if err != nil {
		log.Println("Error deleting ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Deleted ticket with id=%d\n", id)
	return ticket, nil
}
