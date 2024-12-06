package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Payment/internal/models"
)

// PayForTicket processes a payment request by decoding the JSON body,
// validating the card information and ensuring sufficient funds are available.
// It responds with appropriate HTTP status codes based on the outcome of the payment processing.
func PayForTicket(w http.ResponseWriter, r *http.Request) {
	log.Println("Payment process...")

	var paymentData models.PayInfo
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	log.Println("Payment servise got:", string(bs))

	err = json.Unmarshal(bs, &paymentData)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if paymentData.CardInfo.Number == "" {
		log.Println("Missing card number")
		http.Error(w, "Card number is required", http.StatusBadRequest)
		return
	}

	if paymentData.Cost > int(paymentData.CardInfo.Balance) {
		log.Printf("insufficient funds, card: %s, balance: %d, cost: %d\n",
			paymentData.CardInfo.Number, paymentData.CardInfo.Balance, paymentData.Cost)
		http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
		return
	}

	// тут какая-то логика по отправке

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Printf("Payment processed successfully for card: %s\n", paymentData.CardInfo.Number)
}
