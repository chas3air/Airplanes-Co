package router

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Management_tickets/internal/config"
	"github.com/gorilla/mux"
)

var database_url = config.DATABASE_URL

// GetAllTicketsHandler handles a GET request to fetch all tickets.
// It retrieves all tickets from the database and returns them in JSON format.
func GetAllTicketsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all tickets")

	resp, err := http.Get(database_url + "/get")
	if err != nil {
		log.Println("Cannot send request to", database_url)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Bad response: cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully fetched all tickets.")
}

// GetTicketByIdHandler handles a GET request to retrieve a ticket by its ID.
// It fetches the ticket from the database and returns it in JSON format.
func GetTicketByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Retrieving ticket by ID")

	id_s := mux.Vars(r)["id"]

	resp, err := http.Get(database_url + "/get/" + id_s)
	if err != nil {
		log.Println("Cannot send request to", database_url)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully retrieved ticket by ID.")
}

// InsertTicketHandler handles a POST request to add a new ticket.
// It takes ticket data in JSON format from the request body and inserts it into the database.
func InsertTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting ticket")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	resp, err := http.Post(database_url+"/insert", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error sending request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully inserted ticket.")
}

// UpdateTicketHandler handles a PATCH request to update ticket information.
// It takes the updated ticket data in JSON format from the request body.
func UpdateTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating ticket")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req, err := http.NewRequest(http.MethodPatch, database_url+"/update", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error creating request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
	log.Println("Successfully updated ticket.")
}

// DeleteTicketHandler handles a DELETE request to remove a ticket by its ID.
// It deletes the specified ticket from the database.
func DeleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting ticket")

	id_s := mux.Vars(r)["id"]

	req, err := http.NewRequest(http.MethodDelete, database_url+"/delete/"+id_s, nil)
	if err != nil {
		log.Println("Error creating request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(responseBody)
	log.Println("Successfully deleted ticket.")
}
