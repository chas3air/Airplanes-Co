package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Ticket struct {
	Id             uuid.UUID `json:"id" bson:"id"`
	FlightId       uuid.UUID `json:"flightId" bson:"flightId"`
	OwnerId        uuid.UUID `json:"ownerId" bson:"ownerId"`
	TicketCost     float64   `json:"ticketCost" bson:"ticketCost"`
	ClassOfService string    `json:"classOfService" bson:"classOfService"`
}

func (t Ticket) String() string {
	return fmt.Sprintf("Ticket ID: %d, flightId: %d, Owner ID: %d, Cost: %.2f",
		t.Id, t.FlightId, t.OwnerId, t.TicketCost)
}
