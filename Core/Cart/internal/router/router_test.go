package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chas3air/Airplanes-Co/Core/Cart/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestGetTicketsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/cart", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTicketsHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetTicketsHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var tickets map[string][]models.Ticket
	if err := json.Unmarshal(rr.Body.Bytes(), &tickets); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(tickets) != 0 {
		t.Errorf("Expected empty ticket map, got %v", tickets)
	}
}

func TestInsertTicketHandler(t *testing.T) {
	ticket := models.Ticket{
		Id:             uuid.New(),
		FlightInfo:     models.Flight{},
		Owner:          models.Customer{},
		TicketCost:     150,
		ClassOfService: "Economy",
	}
	body, _ := json.Marshal(struct {
		Id     string        `json:"id"`
		Ticket models.Ticket `json:"ticket"`
	}{
		Id:     ticket.Id.String(),
		Ticket: ticket,
	})

	req, err := http.NewRequest("POST", "/cart", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InsertTicketHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("InsertTicketHandler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Проверяем, что билет был добавлен
	if len(TicketsCart) != 1 || len(TicketsCart[ticket.Id.String()]) != 1 || TicketsCart[ticket.Id.String()][0].Id != ticket.Id {
		t.Errorf("Expected ticket not found in cart: got %v", TicketsCart)
	}
}

func TestUpdateTicketHandler(t *testing.T) {
	initialTicket := models.Ticket{
		Id:             uuid.New(),
		FlightInfo:     models.Flight{},
		Owner:          models.Customer{},
		TicketCost:     150,
		ClassOfService: "Economy",
	}

	TicketsCart[initialTicket.Id.String()] = []models.Ticket{initialTicket}

	updatedTicket := initialTicket
	updatedTicket.TicketCost = 200.00
	body, err := json.Marshal(struct {
		Id     string        `json:"id"`
		Ticket models.Ticket `json:"ticket"`
	}{
		Id:     updatedTicket.Id.String(),
		Ticket: updatedTicket,
	})
	if err != nil {
		t.Fatal("Error marshaling updated ticket:", err)
	}

	req, err := http.NewRequest("PATCH", "/cart/"+updatedTicket.Id.String(), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateTicketHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("UpdateTicketHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if TicketsCart[updatedTicket.Id.String()][0].TicketCost != 200.00 {
		t.Errorf("Expected ticket cost to be updated: got %v", TicketsCart[updatedTicket.Id.String()][0].TicketCost)
	}
}

func TestDeleteTicketHandler(t *testing.T) {
	TicketsCart = make(map[string][]models.Ticket)

	generatedId := uuid.New()
	ticket := models.Ticket{
		Id:             generatedId,
		FlightInfo:     models.Flight{},
		Owner:          models.Customer{},
		TicketCost:     150,
		ClassOfService: "Economy",
	}
	TicketsCart[generatedId.String()] = []models.Ticket{ticket}

	req, err := http.NewRequest(http.MethodDelete, "/cart/"+generatedId.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/cart/{id}", DeleteTicketHandler).Methods(http.MethodDelete)

	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DeleteTicketHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if _, ok := TicketsCart[generatedId.String()]; ok {
		t.Errorf("Expected ticket to be deleted, got %v", TicketsCart[generatedId.String()])
	}
}
