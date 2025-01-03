package routes

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/DAL_tickets/internal/models"
	"github.com/chas3air/Airplanes-Co/Core/DAL_tickets/internal/service"
	"github.com/chas3air/Airplanes-Co/Core/DAL_tickets/internal/storage"
	"github.com/gorilla/mux"
)

var TicketsDB = storage.MustGetInstanceOfTicketsStorage("psql")
var limitTime = service.GetLimitTime("PSQL_LIMIT_RESPONSE_TIME")

func GetTickets(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all tickets")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	entities, err := TicketsDB.GetAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Error retrieving tickets")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tickets, ok := entities.([]models.Ticket)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(tickets)
	if err != nil {
		log.Println("Cannot marshal object")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched tickets.")
}

func GetTicketById(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching ticket by ID")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	id := mux.Vars(r)["id"]

	entity, err := TicketsDB.GetAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Error retrieving ticket by ID:", id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tickets, ok := entity.([]models.Ticket)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	var ticket models.Ticket
	for _, t := range tickets {
		if t.Id.String() == id {
			ticket = t
		}
	}

	bs, err := json.Marshal(ticket)
	if err != nil {
		log.Println("Cannot marshal object")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched ticket by ID.")
}

func InsertTicket(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting ticket")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Println("Text comes to dal_tickets:", string(bs))

	var tickets []models.Ticket
	err = json.Unmarshal(bs, &tickets)
	if err != nil {
		log.Println("Cannot unmarshal request body to ticket")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, ticket := range tickets {
		_, err := TicketsDB.Insert(ctx, ticket)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				log.Println("Request timed out")
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
				return
			} else {
				log.Printf("Error inserting ticket with ID: %d\n", ticket.Id)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	log.Println("Successfully inserted ticket")
}

func UpdateTicket(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating ticket")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var ticket models.Ticket
	err = json.Unmarshal(bs, &ticket)
	if err != nil {
		log.Println("Cannot unmarshal request body to ticket")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := TicketsDB.Update(ctx, ticket)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error updating ticket with ID: %d\n", ticket.Id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal updated ticket")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully updated ticket with ID: %d", ticket.Id)
}

func DeleteTicket(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting ticket")

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	id := mux.Vars(r)["id"]

	obj, err := TicketsDB.Delete(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Println("Error deleting ticket with ID:", id)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal deleted ticket")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully deleted ticket with ID: %v", id)
}
