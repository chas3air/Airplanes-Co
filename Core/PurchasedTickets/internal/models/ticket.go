package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Ticket struct {
	Id             uuid.UUID `json:"id" bson:"id"`
	FlightInfo     Flight    `json:"flightInfo" bson:"flightInfo"`
	Owner          Customer  `json:"owner" bson:"owner"`
	TicketCost     int   `json:"ticketCost" bson:"ticketCost"`
	ClassOfService string    `json:"classOfService" bson:"classOfService"`
}

func (t Ticket) String() string {
	return fmt.Sprintf("Ticket ID: %d, flightInfo: %v, Owner: %v, Cost: %d",
		t.Id, t.FlightInfo, t.Owner, t.TicketCost)
}
