package cart

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/chas3air/Airplanes-Co/Cart/internal/models"
	"github.com/gorilla/mux"
)

var TicketsCart = make([]models.Ticket, 0, 10)

// GetAllTicketsHandler handles GET requests to fetch all tickets in the cart.
// It returns the tickets as a JSON array.
func GetAllTicketsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching tickets from cart")

	bs, err := json.Marshal(TicketsCart)
	if err != nil {
		log.Println("Error marshaling tickets:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched tickets.")
}

// InsertTicketHandler handles POST requests to insert a new ticket into the cart.
// It reads the ticket data from the request body and adds it to the cart.
func InsertTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting ticket to cart")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var ticket models.Ticket
	err = json.Unmarshal(body, &ticket)
	if err != nil {
		log.Println("Error unmarshaling ticket:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	TicketsCart = append(TicketsCart, ticket)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully inserted ticket.")
}

// UpdateTicketHandler handles PATCH requests to update an existing ticket in the cart.
// It reads the updated ticket data from the request body and modifies the existing entry.
func UpdateTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating ticket in cart")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var ticket models.Ticket
	err = json.Unmarshal(body, &ticket)
	if err != nil {
		log.Println("Error unmarshaling ticket:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update existing ticket in the cart
	for i, v := range TicketsCart {
		if v.Id == ticket.Id {
			TicketsCart[i] = ticket
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully updated ticket.")
}

// DeleteTicketHandler handles DELETE requests to remove a ticket from the cart by its ID.
// It returns the deleted ticket as a JSON object.
func DeleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting ticket from cart")

	id_s, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot get id from URL")
		http.Error(w, "Bad request: cannot get id from URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Bad request: invalid ticket ID:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var deletedTicket models.Ticket
	var newCart []models.Ticket
	for _, v := range TicketsCart {
		if v.Id != id {
			newCart = append(newCart, v)
		} else {
			deletedTicket = v
		}
	}

	TicketsCart = newCart

	bs, err := json.Marshal(deletedTicket)
	if err != nil {
		log.Println("Error marshaling deleted ticket:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully deleted ticket.")
}

// ClearHandler handles requests to clear all tickets from the cart.
// It resets the cart to an empty state.
func ClearHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Clearing cart")

	TicketsCart = make([]models.Ticket, 0, 10)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Println("Successfully cleared cart.")
}
