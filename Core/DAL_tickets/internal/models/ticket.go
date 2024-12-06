package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Ticket struct {
	Id             uuid.UUID `json:"id" bson:"id"`
	FlightInfo     Flight    `json:"flightInfo" bson:"flightInfo"`
	Owner          Customer  `json:"owner" bson:"owner"`
	ClassId        int       `json:"classId" bson:"classId"`
	TicketCost     int       `json:"ticketCost" bson:"ticketCost"`
	ClassOfService ClassName `json:"classOfService" bson:"classOfService"`
}

type ClassName struct {
	Id    int    `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
}

func (t Ticket) String() string {
	return fmt.Sprintf("Ticket ID: %d, flightInfo: %v, Owner: %v, Cost: %d",
		t.Id, t.FlightInfo, t.Owner, t.TicketCost)
}
