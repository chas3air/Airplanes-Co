package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Cart/internal/models"
	"github.com/gorilla/mux"
)

var TicketsCart = make([]models.Ticket, 0, 10)

// GetAllTicketsHandler handles GET requests to fetch all tickets in the cart.
// It returns the tickets as a JSON array.
func GetTicketsHandler(w http.ResponseWriter, r *http.Request) {
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

	var isExist bool = false
	for i, v := range TicketsCart {
		if v.Id == ticket.Id {
			isExist = true
			TicketsCart[i] = ticket
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if !isExist {
		w.WriteHeader(http.StatusNoContent)
		log.Println("Ticket not found.")
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
		log.Println("Successfully updated ticket.")
	}

}

// DeleteTicketHandler handles DELETE requests to remove a ticket from the cart by its ID.
// It returns the deleted ticket as a JSON object.
func DeleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting ticket from cart")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Printf("Bad request: cannot get id from URL, id: %v\n", id)
		http.Error(w, "Bad request: cannot get id from URL", http.StatusBadRequest)
		return
	}

	var deletedTicket models.Ticket
	var newCart []models.Ticket
	found := false

	for _, v := range TicketsCart {
		if v.Id.String() != id {
			newCart = append(newCart, v)
		} else {
			deletedTicket = v
			found = true
		}
	}

	if found {
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
	} else {
		log.Println("Ticket not found:", id)
		http.Error(w, "Ticket not found", http.StatusNotFound)
	}
}

// ClearHandler handles requests to clear all tickets from the cart.
// It resets the cart to an empty state and returns all tickets that were in the cart.
func ClearHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Clearing cart")

	ticketsToReturn := TicketsCart

	TicketsCart = TicketsCart[:0]

	bs, err := json.Marshal(ticketsToReturn)
	if err != nil {
		log.Println("Cannot marshal tickets for response:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully cleared cart and returned tickets.")
}
