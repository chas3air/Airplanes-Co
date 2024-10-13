package models

import (
	"fmt"
	"time"
)

type Ticket struct {
	Id                int       `json:"id" bson:"id" gorm:"primaryKey"`
	FlightFromWhere   string    `json:"flightFromWhere" bson:"flightFromWhere" gorm:"column:flightFromWhere;not null"`
	FlightDestination string    `json:"flightDestination" bson:"flightDestination" gorm:"column:flightDestination;not null"`
	FlightTime        time.Time `json:"flightTime" bson:"flightTime" gorm:"column:flightTime;not null"`
	OwnerId           int       `json:"ownerId" bson:"ownerId" gorm:"column:ownerId;not null"`
	TicketCost        float64   `json:"ticketCost" bson:"ticketCost" gorm:"column:ticketCost;not null"`
	ClassOfService    string    `json:"classOfService" bson:"classOfService" gorm:"column:classOfService;not null"`
}

func (t Ticket) String() string {
	return fmt.Sprintf("Ticket ID: %d, From: %s, To: %s, Flight Time: %s, Owner ID: %d, Cost: %.2f",
		t.Id, t.FlightFromWhere, t.FlightDestination, t.FlightTime.Format(time.RFC3339), t.OwnerId, t.TicketCost)
}
