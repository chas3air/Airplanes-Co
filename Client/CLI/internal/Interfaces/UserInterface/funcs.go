package ui

import (
	"errors"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
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

func ClearCart(userId uuid.UUID) error {
	req, err := http.NewRequest(http.MethodDelete, config.Backend_url+"/cart/clear/"+userId.String(), nil)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Timeout: service.GetLimitTime(),
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return errors.New("Response status: " + resp.Status)
	}

	return nil
}
