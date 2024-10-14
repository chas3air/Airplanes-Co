package models

import (
	"fmt"
)

type Ticket struct {
	Id             int     `json:"id" bson:"id"`
	FlightId       int     `json:"flightId" bson:"flightId"`
	OwnerId        int     `json:"ownerId" bson:"ownerId"`
	TicketCost     float64 `json:"ticketCost" bson:"ticketCost"`
	ClassOfService string  `json:"classOfService" bson:"classOfService"`
}

func (t Ticket) String() string {
	return fmt.Sprintf("Ticket ID: %d, flightId: %d, Owner ID: %d, Cost: %.2f",
		t.Id, t.FlightId, t.OwnerId, t.TicketCost)
}
