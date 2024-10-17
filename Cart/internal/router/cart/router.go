package cart

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/chas3air/Airplanes-Co/Cart/internal/config"
	"github.com/chas3air/Airplanes-Co/Cart/internal/models"
	"github.com/gorilla/mux"
)

func GetAllTicketsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching tickets from cart")

	bs, err := json.Marshal(config.TicketsCart)
	if err != nil {
		log.Println("Cannot marshal tockets")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully fetched tickets.")
}

func InsertTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Inserting tickets to cart")

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
		log.Println("Cannot unmarshal ticket")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	config.TicketsCart = append(config.TicketsCart, ticket)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully inserted ticket.")
}

func UpdateTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating ticket")

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
		log.Println("Cannot unmarshal ticket")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Изменение существующего билета в корзине
	for i, v := range config.TicketsCart {
		if v.Id == ticket.Id {
			config.TicketsCart[i] = ticket
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully updating ticket.")
}

func DeleteTicketHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Deleting ticket")

	id_s, ok := mux.Vars(r)["id"]
	if !ok {
		log.Println("Bad request: cannot get id from url string")
		http.Error(w, "Bad request: cannot get id from url string", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Bad request: wrong ticket ID")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var deletedTicket models.Ticket
	var newCart = make([]models.Ticket, 0, len(config.TicketsCart))
	for _, v := range config.TicketsCart {
		if v.Id != id {
			newCart = append(newCart, v)
		} else {
			deletedTicket = v
		}
	}

	copy(config.TicketsCart, newCart)
	bs, err := json.Marshal(deletedTicket)
	if err != nil {
		log.Println("Cannot marshal ticket")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Successfully deleting ticket.")
}
