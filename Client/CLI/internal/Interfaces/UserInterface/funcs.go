package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
)

func CreateTicket(flight models.Flight) (models.Ticket, error) {
	var isValid bool
	var indexOfClass int
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println(flight.FlightSeatsCosts)
	fmt.Println(config.NamesSeats)
	fmt.Println("Enter cost for seat from list")
	_ = scanner.Scan()
	cost := scanner.Text()

	for i, v := range flight.FlightSeatsCosts {
		if strconv.Itoa(v) == cost {
			indexOfClass = i
			isValid = true
		}
	}
	if !isValid {
		return models.Ticket{}, fmt.Errorf("underfined flights seat cost")
	}

	ticket := models.Ticket{
		FlightInfo:     flight,
		TicketCost:     flight.FlightSeatsCosts[indexOfClass],
		ClassOfService: config.NamesSeats[indexOfClass],
	}

	return ticket, nil
}

func BuyTicket(ticket models.Ticket) (models.Ticket, error) {

	return models.Ticket{}, nil
}
