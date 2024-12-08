package models

import (
	"time"

	"github.com/google/uuid"
)

type Flight struct {
	Id               uuid.UUID `json:"id" bson:"id"`
	FromWhere        string    `json:"fromWhere" bson:"fromWhere"`
	Destination      string    `json:"destination" bson:"destination"`
	FlightTime       time.Time `json:"flightTime" bson:"flightTime"`
	FlightDuration   int       `json:"flightDuration" bson:"flightDuration"`
	FlightSeatsCosts []int     `json:"flightSeatsCost" bson:"flightSeatsCost"`
	Airplane         uuid.UUID `json:"airplane"`
}
