package router

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/OrderTicket/internal/models"
	"github.com/chas3air/Airplanes-Co/OrderTicket/internal/service"
	"github.com/google/uuid"
)

var cart_url = os.Getenv("CART_URL")
var limitTime = service.GetLimitTime()

var httpClient = &http.Client{
	Timeout: limitTime,
}

// OrderTicketHandler processes ticket ordering requests by reading ticket data from the request body,
// assigning a new unique ID, and sending the data to the cart service for processing.
// It handles errors and returns appropriate HTTP responses.
func OrderTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Order ticket processing...")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var ticket models.Ticket
	err = json.Unmarshal(body, &ticket)
	if err != nil {
		log.Println("Cannot unmarshal request body to ticket")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ticket.Id = uuid.New()
	ticket_bs, err := json.Marshal(ticket)
	if err != nil {
		log.Println("Cannot marshal ticket to JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := httpClient.Post(cart_url, "application/json", bytes.NewBuffer(ticket_bs)) // Используем httpClient
	if err != nil {
		log.Printf("Error sending request to %s: %v", cart_url, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	response_body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response_body)
	log.Println("Successfully ordered ticket.")
}
