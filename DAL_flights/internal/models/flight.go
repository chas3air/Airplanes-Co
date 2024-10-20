package models

import (
	"fmt"
	"time"
)

type Flight struct {
	Id             int       `json:"id" bson:"id"`
	FromWhere      string    `json:"fromWhere" bson:"fromWhere"`
	Destination    string    `json:"destination" bson:"destination"`
	FlightTime     time.Time `json:"flightTime" bson:"flightTime"`
	FlightDuration int       `json:"flightDuration" bson:"flightDuration"`
}

func (f Flight) String() string {
	return fmt.Sprintf("Flight ID: %d, From: %s, To: %s, Departure: %s, Duration: %d minutes",
		f.Id, f.FromWhere, f.Destination, f.FlightTime.Format(time.RFC3339), f.FlightDuration)
}
