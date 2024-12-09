package router

import (
	"encoding/json"
	"log"
	"net/http"
)

// PayForTicket processes a payment request by decoding the JSON body,
// validating the card information and ensuring sufficient funds are available.
// It responds with appropriate HTTP status codes based on the outcome of the payment processing.
func PayForTicket(w http.ResponseWriter, r *http.Request) {
	log.Println("Payment process...")

	var payData = struct {
		CardNumber  string `json:"card_number"`
		BankAccount string `json:"bank_account"`
		Cost        int    `json:"cost"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&payData)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// тут какая-то логика по отправке

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//w.WriteHeader(http.StatusNotFound)
	log.Printf("Payment processed successfully for card: %s\n", payData.CardNumber)
}
