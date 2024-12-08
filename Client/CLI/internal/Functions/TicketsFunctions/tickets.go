package ticketsfunctions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/config"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/models"
	"github.com/chas3air/Airplanes-Co/Client/CLI/internal/service"
	"github.com/google/uuid"
)

var limitTime = service.GetLimitTime()

func SendTicketToTheCart(ticket models.Ticket) (models.Ticket, error) {
	bs, err := json.Marshal(ticket)
	if err != nil {
		return models.Ticket{}, fmt.Errorf("error marshalling ticket: %w", err)
	}

	httpClient := &http.Client{
		Timeout: limitTime,
	}

	resp, err := httpClient.Post(config.Backend_url+"/cart"+"?ownerId="+ticket.Owner.Id.String(), "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return models.Ticket{}, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return models.Ticket{}, fmt.Errorf("received response: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&ticket)
	if err != nil {
		return models.Ticket{}, fmt.Errorf("error decoding response: %w", err)
	}

	return ticket, nil
}

func GetTicketFromCart(userId uuid.UUID) ([]models.Ticket, error) {
	var tickets []models.Ticket
	httpClient := &http.Client{
		Timeout: limitTime,
	}

	resp, err := httpClient.Get(config.Backend_url + "/cart?ownerId=" + userId.String())
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("receiver response: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&tickets)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return tickets, nil
}

func GetPurchasedTickets(user models.Customer) ([]models.Ticket, error) {
	var purchasedTickets []models.Ticket

	httpClient := &http.Client{
		Timeout: limitTime,
	}

	resp, err := httpClient.Get(config.Backend_url + "/purchasedTickets?ownerId=" + url.PathEscape(user.Id.String()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("received response: %s", resp.Status)
	}

	err = json.NewDecoder(resp.Body).Decode(&purchasedTickets)
	if err != nil {
		return nil, err
	}

	return purchasedTickets, nil
}

func PrintTickets(tickets []models.Ticket, exclude ...string) {
	excludeMap := make(map[string]struct{})
	for _, col := range exclude {
		excludeMap[col] = struct{}{}
	}

	headers := []string{"ID", "Flight Info", "Owner", "Ticket Cost", "Class"}
	widths := []int{36, 50, 50, 15, 20}
	var selectedHeaders []string
	var selectedWidths []int

	for i, header := range headers {
		if _, ok := excludeMap[header]; !ok {
			selectedHeaders = append(selectedHeaders, header)
			selectedWidths = append(selectedWidths, widths[i])
		}
	}

	for i, header := range selectedHeaders {
		fmt.Printf("| %-*s ", selectedWidths[i], header)
	}
	fmt.Println("|")
	fmt.Println(strings.Repeat("-", getTotalWidth(selectedWidths)))

	for _, ticket := range tickets {
		if _, ok := excludeMap["ID"]; !ok {
			fmt.Printf("| %-36s ", ticket.Id.String())
		}
		if _, ok := excludeMap["Flight Info"]; !ok {
			flightInfo := fmt.Sprintf("%s to %s", ticket.FlightInfo.FromWhere, ticket.FlightInfo.Destination)
			fmt.Printf("| %-50s ", flightInfo)
		}
		if _, ok := excludeMap["Owner"]; !ok {
			ownerInfo := fmt.Sprintf("%s %s", ticket.Owner.Name, ticket.Owner.Surname)
			fmt.Printf("| %-50s ", ownerInfo)
		}
		if _, ok := excludeMap["Ticket Cost"]; !ok {
			fmt.Printf("| %-15d ", ticket.TicketCost)
		}
		if _, ok := excludeMap["Class"]; !ok {
			fmt.Printf("| %-20s ", ticket.ClassOfService)
		}
		fmt.Println("|")
	}

	fmt.Println(strings.Repeat("-", getTotalWidth(selectedWidths)))
}

func getTotalWidth(widths []int) int {
	total := 0
	for _, w := range widths {
		total += w + 3
	}
	return total
}

func PayForTickets(card_number string) error {
	bs, err := json.Marshal(card_number)
	if err != nil {
		return err
	}

	httpClient := &http.Client{
		Timeout: limitTime,
	}

	resp, err := httpClient.Post(config.Backend_url+"/payment", "application/json", bytes.NewReader(bs))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("response code: %d", resp.StatusCode)
	}

	return nil
}
