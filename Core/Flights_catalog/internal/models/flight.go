package models

import (
	"fmt"
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

func (f Flight) String() string {
	return fmt.Sprintf("Flight ID: %s, From: %s, To: %s, Departure: %s, Duration: %d minutes, Flighys cost seats: %v",
		f.Id.String(), f.FromWhere, f.Destination, f.FlightTime.Format(time.RFC3339), f.FlightDuration, f.FlightSeatsCosts)
}
