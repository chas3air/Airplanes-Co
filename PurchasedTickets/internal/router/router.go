package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/PurchasedTickets/internal/models"
	"github.com/chas3air/Airplanes-Co/PurchasedTickets/internal/service"
)

var managementTicketsURL = os.Getenv("MANAGEMENT_TICKETS_URL")
var limitTime = service.GetLimitTime()

// GetPurchasedTicketsHandler handles requests to retrieve purchased tickets.
// It extracts the owner ID from the request parameters,
// sends a GET request to an external service to fetch the list of tickets,
// filters the tickets by owner ID, and returns the filtered list in JSON format.
func GetPurchasedTicketsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Processing request to get purchased tickets...")

	ownerId := r.FormValue("ownerId")

	resp, err := http.Get(managementTicketsURL)
	if err != nil {
		log.Printf("Error sending GET request to %s: %v", managementTicketsURL, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-200 response from management service: %s", resp.Status)
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tickets []models.Ticket
	if err := json.Unmarshal(bs, &tickets); err != nil {
		log.Println("Error unmarshalling response body to tickets:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var outTickets []models.Ticket
	for _, ticket := range tickets {
		if ticket.OwnerId.String() == ownerId {
			outTickets = append(outTickets, ticket)
		}
	}

	out, err := json.Marshal(outTickets)
	if err != nil {
		log.Println("Error marshalling output tickets to JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
	log.Println("Successfully retrieved purchased tickets.")
}
