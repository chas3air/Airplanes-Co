package ui

import (
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/google/uuid"
)

func CreateTicket(userID uuid.UUID, flight models.Flight, cost int, classOfService string) models.Ticket {
	ticket := models.Ticket{
		Id:         uuid.New(),
		FlightInfo: flight,
		Owner: models.Customer{
			Id: userID,
		},
		TicketCost:     cost,
		ClassOfService: classOfService,
	}

	return ticket
}