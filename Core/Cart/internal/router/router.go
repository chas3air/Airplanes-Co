package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Cart/internal/models"
	"github.com/gorilla/mux"
)

// TicketsCart - a map that stores tickets, where the key is a string (ID) and the value is a slice of tickets.
var TicketsCart = make(map[string][]models.Ticket)

// writeJSONResponse writes a JSON response and sets the Content-Type header.
// Parameters:
//   - w: http.ResponseWriter for sending the response to the client.
//   - status: HTTP status code for the response.
//   - data: data to be sent as JSON.
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Println("Error writing JSON response:", err)
		}
	}
}

// readRequestBody reads the request body and returns it as []byte.
// Parameter:
//   - r: http.Request representing the incoming request.
//
// Returns:
//   - []byte: the content of the request body.
//   - error: an error if there was a problem reading the body.
func readRequestBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return body, nil
}

// GetTicketsHandler handles GET requests to fetch all tickets in the cart.
// Parameters:
//   - w: http.ResponseWriter for sending the response to the client.
//   - r: http.Request representing the incoming request.
//
// Response:
//   - Returns JSON containing the current tickets in the cart and a status of 200 OK.
func GetTicketsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching tickets from cart")
	id := mux.Vars(r)["id"]
	tickets := TicketsCart[id]
	writeJSONResponse(w, http.StatusOK, tickets)
}

// InsertTicketHandler handles POST requests to insert a new ticket into the cart.
// Parameters:
//   - w: http.ResponseWriter for sending the response to the client.
//   - r: http.Request representing the incoming request.
//
// Response:
//   - Returns a status of 201 Created and the added ticket in JSON format,
//     or a status of 400 Bad Request in case of an error.
func InsertTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting ticket to cart")

	var ticket models.Ticket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		log.Println("Cannot marshal request body to object")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Ticket:", ticket.String())

	// Add the ticket to the cart
	TicketsCart[ticket.Owner.Id.String()] = append(TicketsCart[ticket.Owner.Id.String()], ticket)
	log.Println("Successfully inserted ticket.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	log.Println("Item successfully created")
}

// UpdateTicketHandler handles PATCH requests to update an existing ticket in the cart.
// Parameters:
//   - w: http.ResponseWriter for sending the response to the client.
//   - r: http.Request representing the incoming request.
//
// Response:
//   - Returns a status of 200 OK and the updated ticket in JSON format,
//     or a status of 404 Not Found if the ticket is not found.
func UpdateTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating ticket in cart")

	var ticket models.Ticket
	err := json.NewDecoder(r.Body).Decode(&ticket)
	if err != nil {
		log.Println("Cannot marshal request body to object")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add the ticket to the cart
	TicketsCart[ticket.Owner.Id.String()] = append(TicketsCart[ticket.Owner.Id.String()], ticket)
	log.Println("Successfully inserted ticket.")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	log.Println("Item successfully updated")
}

// DeleteTicketHandler handles DELETE requests to remove a ticket from the cart by its ID.
// Parameters:
//   - w: http.ResponseWriter for sending the response to the client.
//   - r: http.Request representing the incoming request.
//
// Response:
//   - Returns a status of 200 OK and a message that the ticket was deleted,
//     or a status of 404 Not Found if the ticket is not found.
func DeleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting ticket from cart")

	id, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot get id from URL")
		writeJSONResponse(w, http.StatusBadRequest, "Bad request: cannot get id from URL")
		return
	}

	log.Printf("Deleting ticket with ID: %s\n", id)

	if _, ok := TicketsCart[id]; ok {
		delete(TicketsCart, id)
		writeJSONResponse(w, http.StatusOK, "Ticket deleted")
		log.Println("Successfully deleted ticket.")
	} else {
		log.Println("Ticket not found for deletion.")
		writeJSONResponse(w, http.StatusNotFound, "Ticket not found")
	}
}

// ClearHandler handles requests to clear all tickets from the cart.
// Parameters:
//   - w: http.ResponseWriter for sending the response to the client.
//   - r: http.Request representing the incoming request.
//
// Response:
//   - Returns a status of 200 OK and the list of all tickets that were in the cart before clearing.
func ClearHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Clearing cart")
	id := mux.Vars(r)["id"]

	var tickets = TicketsCart[id]
	delete(TicketsCart, id)

	log.Println("Size of ticket card of current user:", len(TicketsCart[id]))

	writeJSONResponse(w, http.StatusOK, tickets)
}
