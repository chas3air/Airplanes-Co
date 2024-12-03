package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chas3air/Airplanes-Co/Core/DAL_tickets/internal/config"
	"github.com/chas3air/Airplanes-Co/Core/DAL_tickets/internal/models"
)

type PsqlTicketsStorage struct {
	PsqlStorage
}

// MustNewPsqlTicketsStorage initializes a new PsqlTicketsStorage instance.
// It checks the database connection and panics if the tickets table is unavailable.
func MustNewPsqlTicketsStorage() PsqlTicketsStorage {
	return PsqlTicketsStorage{
		PsqlStorage: PsqlStorage{},
	}
}

// GetAll retrieves all tickets from the database.
func (s PsqlTicketsStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.GetAll"

	// Открываем соединение
	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, `
		SELECT * FROM `+os.Getenv("PSQL_TABLE_NAME")+`;
	`)
	if err != nil {
		log.Println("Error querying tickets:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	var ticket models.Ticket
	tickets := make([]models.Ticket, 0, 10)

	for rows.Next() {
		err := rows.Scan(&ticket.Id, &ticket.FlightInfo.Id, &ticket.Owner.Id,
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

	log.Printf("Retrieved %d tickets\n", len(tickets))
	return tickets, nil
}

// GetById retrieves a ticket by its ID.
func (s PsqlTicketsStorage) GetById(ctx context.Context, id any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.GetById"

	// Открываем соединение
	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	row := db.QueryRowContext(ctx, `
		SELECT * FROM `+os.Getenv("PSQL_TABLE_NAME")+`
		WHERE id = $1;
	`, id)

	var ticket models.Ticket

	err = row.Scan(&ticket.Id, &ticket.FlightInfo.Id, &ticket.Owner.Id, &ticket.TicketCost, &ticket.ClassOfService)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No ticket found with ID=%s\n", id)
			return nil, fmt.Errorf("%s: no ticket found with ID=%d", op, id)
		}
		log.Println("Error scanning ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrieved ticket with ID=%v\n", ticket.Id)
	return ticket, nil
}

// Insert adds a new ticket to the database.
func (s PsqlTicketsStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Insert"

	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	ticket := innerObj.(models.Ticket)

	_, err = db.ExecContext(ctx, `
		INSERT INTO `+os.Getenv("PSQL_TABLE_NAME")+`
		(id, flightId, ownerId, ticketCost, classOfService)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`, ticket.Id, ticket.FlightInfo.Id, ticket.Owner.Id, ticket.TicketCost, ticket.ClassOfService)

	if err != nil {
		log.Println("Error inserting ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Inserted ticket with ID=%s\n", ticket.Id.String())
	return ticket, nil
}

// Update modifies an existing ticket in the database.
func (s PsqlTicketsStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Update"

	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	ticket := innerObj.(models.Ticket)

	_, err = db.ExecContext(ctx, `
		UPDATE `+os.Getenv("PSQL_TABLE_NAME")+`
		SET flightId = $1, ownerId = $2, ticketCost = $3, classOfService = $4
		WHERE id = $5;
	`, ticket.FlightInfo.Id, ticket.Owner.Id, ticket.TicketCost, ticket.ClassOfService, ticket.Id)

	if err != nil {
		log.Println("Error updating ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Updated ticket with ID=%d\n", ticket.Id)
	return ticket, nil
}

// Delete removes a ticket from the database by its ID.
func (s PsqlTicketsStorage) Delete(ctx context.Context, id any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Delete"

	// Открываем соединение
	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	ticket, err := s.GetById(ctx, id)
	if err != nil {
		log.Println("Error getting ticket by ID:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	_, err = db.ExecContext(ctx, `
		DELETE FROM `+os.Getenv("PSQL_TABLE_NAME")+`
		WHERE id = $1;
	`, id)

	if err != nil {
		log.Println("Error deleting ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Deleted ticket with ID=%s\n", id)
	return ticket, nil
}
