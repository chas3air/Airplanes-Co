package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chas3air/Airplanes-Co/DAL_tickets/internal/models"
)

type PsqlTicketsStorage struct {
	PsqlStorage
}

// MustNewPsqlTicketsStorage initializes a new PsqlTicketsStorage instance.
// It checks the database connection and panics if the tickets table is unavailable.
func MustNewPsqlTicketsStorage(db *sql.DB) PsqlTicketsStorage {
	const op = "DAL.internal.storage.psqlRepository.newTicketStorage"
	err := db.Ping()
	if err != nil {
		log.Println("Tickets table is unavailable: " + err.Error())
		panic(fmt.Errorf("%s: %w", op, err))
	}

	return PsqlTicketsStorage{
		PsqlStorage: NewPsqlStorage(db),
	}
}

// GetAll retrieves all tickets from the database.
// It returns a slice of Ticket models and an error, if any occurs.
func (s PsqlTicketsStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.GetAll"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Tickets table is unavailable: " + err.Error())
		return nil, fmt.Errorf("tickets table is unavailable, file: %s: %w", op, err)
	}

	rows, err := s.DB.QueryContext(ctx, `
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

	log.Printf("Retrieved %d tickets\n", len(tickets))
	return tickets, nil
}

// GetById retrieves a ticket by its ID.
// It returns the Ticket model and an error if no ticket is found or if an error occurs.
func (s PsqlTicketsStorage) GetById(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.GetById"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Tickets table is unavailable: " + err.Error())
		return nil, fmt.Errorf("tickets table is unavailable, file: %s: %w", op, err)
	}

	row := s.DB.QueryRowContext(ctx, `
		SELECT * FROM `+os.Getenv("PSQL_TABLE_NAME")+`
		WHERE id = $1;
	`, id)

	var ticket models.Ticket

	err = row.Scan(&ticket.Id, &ticket.FlightId, &ticket.OwnerId, &ticket.TicketCost, &ticket.ClassOfService)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No ticket found with ID=%d\n", id)
			return nil, fmt.Errorf("%s: no ticket found with ID=%d", op, id)
		}
		log.Println("Error scanning ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrieved ticket with ID=%d\n", ticket.Id)
	return ticket, nil
}

// Insert adds a new ticket to the database.
// It returns the inserted Ticket model and an error if the insertion fails.
func (s PsqlTicketsStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Insert"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Tickets table is unavailable: " + err.Error())
		return nil, fmt.Errorf("tickets table is unavailable, file: %s: %w", op, err)
	}

	ticket := innerObj.(models.Ticket)
	var id int

	err = s.DB.QueryRowContext(ctx, `
		INSERT INTO `+os.Getenv("PSQL_TABLE_NAME")+`
		(flightId, ownerId, ticketCost, classOfService)
		VALUES ($1, $2, $3, $4) RETURNING id
	`, ticket.FlightId, ticket.OwnerId, ticket.TicketCost, ticket.ClassOfService).Scan(&id)

	if err != nil {
		log.Println("Error inserting ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Inserted ticket with ID=%d\n", id)
	return s.GetById(ctx, id)
}

// Update modifies an existing ticket in the database.
// It returns the updated Ticket model and an error if the update fails.
func (s PsqlTicketsStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Update"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Tickets table is unavailable: " + err.Error())
		return nil, fmt.Errorf("tickets table is unavailable, file: %s: %w", op, err)
	}

	ticket := innerObj.(models.Ticket)

	_, err = s.DB.ExecContext(ctx, `
		UPDATE `+os.Getenv("PSQL_TABLE_NAME")+`
		SET flightId = $1, ownerId = $2, ticketCost = $3, classOfService = $4
		WHERE id = $5;
	`, ticket.FlightId, ticket.OwnerId, ticket.TicketCost, ticket.ClassOfService, ticket.Id)

	if err != nil {
		log.Println("Error updating ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Updated ticket with ID=%d\n", ticket.Id)
	return s.GetById(ctx, ticket.Id)
}

// Delete removes a ticket from the database by its ID.
// It returns the deleted Ticket model and an error if the deletion fails.
func (s PsqlTicketsStorage) Delete(ctx context.Context, id int) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlTickets.Delete"

	err := s.DB.Ping()
	if err != nil {
		log.Println("Tickets table is unavailable: " + err.Error())
		return nil, fmt.Errorf("tickets table is unavailable, file: %s: %w", op, err)
	}

	ticket, err := s.GetById(ctx, id)
	if err != nil {
		log.Println("Error getting ticket by ID:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	_, err = s.DB.ExecContext(ctx, `
		DELETE FROM `+os.Getenv("PSQL_TABLE_NAME")+`
		WHERE id = $1;
	`, id)

	if err != nil {
		log.Println("Error deleting ticket:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Deleted ticket with ID=%d\n", id)
	return ticket, nil
}
