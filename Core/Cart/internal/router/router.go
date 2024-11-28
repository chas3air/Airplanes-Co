package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Cart/internal/models"
	"github.com/gorilla/mux"
)

var TicketsCart = make([]models.Ticket, 0, 10)

// writeJSONResponse записывает JSON-ответ и устанавливает заголовок Content-Type
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println("Error writing JSON response:", err)
		}
	}
}

// readRequestBody читает тело запроса и возвращает его как []byte
func readRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

// GetTicketsHandler handles GET requests to fetch all tickets in the cart.
func GetTicketsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching tickets from cart")
	writeJSONResponse(w, http.StatusOK, TicketsCart)
}

// InsertTicketHandler handles POST requests to insert a new ticket into the cart.
func InsertTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting ticket to cart")

	body, err := readRequestBody(r)
	if err != nil {
		log.Println("Bad request: cannot read request body:", err)
		writeJSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var ticket models.Ticket
	if err = json.Unmarshal(body, &ticket); err != nil {
		log.Println("Error unmarshaling ticket:", err)
		writeJSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	TicketsCart = append(TicketsCart, ticket)
	writeJSONResponse(w, http.StatusOK, ticket)
	log.Println("Successfully inserted ticket.")
}

// UpdateTicketHandler handles PATCH requests to update an existing ticket in the cart.
func UpdateTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating ticket in cart")

	body, err := readRequestBody(r)
	if err != nil {
		log.Println("Bad request: cannot read request body:", err)
		writeJSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	var ticket models.Ticket
	if err = json.Unmarshal(body, &ticket); err != nil {
		log.Println("Error unmarshaling ticket:", err)
		writeJSONResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	for i, v := range TicketsCart {
		if v.Id == ticket.Id {
			TicketsCart[i] = ticket
			writeJSONResponse(w, http.StatusOK, ticket)
			log.Println("Successfully updated ticket.")
			return
		}
	}

	log.Println("Ticket not found.")
	writeJSONResponse(w, http.StatusNotFound, "Ticket not found")
}

// DeleteTicketHandler handles DELETE requests to remove a ticket from the cart by its ID.
func DeleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting ticket from cart")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot get id from URL")
		writeJSONResponse(w, http.StatusBadRequest, "Bad request: cannot get id from URL")
		return
	}
	log.Printf("Deleting ticket with ID: %s\n", id)

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
		writeJSONResponse(w, http.StatusOK, deletedTicket)
		log.Println("Successfully deleted ticket.")
	} else {
		log.Println("Ticket not found:", id)
		writeJSONResponse(w, http.StatusNotFound, "Ticket not found")
	}
}

// ClearHandler handles requests to clear all tickets from the cart.
func ClearHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Clearing cart")

	ticketsToReturn := TicketsCart
	TicketsCart = []models.Ticket{} // Очистка корзины

	writeJSONResponse(w, http.StatusOK, ticketsToReturn)
	log.Println("Successfully cleared cart and returned tickets.")
}