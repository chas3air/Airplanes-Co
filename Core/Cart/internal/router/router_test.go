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

	var tickets []models.Ticket
	if err := json.Unmarshal(rr.Body.Bytes(), &tickets); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(tickets) != 0 {
		t.Errorf("Expected empty ticket list, got %v", tickets)
	}
}

func TestInsertTicketHandler(t *testing.T) {
	ticket := models.Ticket{
		Id:             uuid.New(),
		FlightInfo:     models.Flight{},
		Owner:          models.Customer{},
		TicketCost:     150.75,
		ClassOfService: "Economy",
	}
	body, _ := json.Marshal(ticket)

	req, err := http.NewRequest("POST", "/cart", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(InsertTicketHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("InsertTicketHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем, что билет был добавлен
	if len(TicketsCart) != 1 || TicketsCart[0].Id != ticket.Id {
		t.Errorf("Expected ticket not found in cart: got %v", TicketsCart)
	}
}

func TestUpdateTicketHandler(t *testing.T) {
	initialTicket := models.Ticket{
		Id:             uuid.New(),
		FlightInfo:     models.Flight{},
		Owner:          models.Customer{},
		TicketCost:     150.75,
		ClassOfService: "Economy",
	}

	TicketsCart = append(TicketsCart, initialTicket)

	updatedTicket := initialTicket
	updatedTicket.TicketCost = 200.00
	TicketsCart[0].TicketCost = 200.00
	body, err := json.Marshal(updatedTicket)
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

	if TicketsCart[0].TicketCost != 200.00 {
		t.Errorf("Expected ticket cost to be updated: got %v", TicketsCart[0].TicketCost)
	}
}

func TestDeleteTicketHandler(t *testing.T) {
	TicketsCart = make([]models.Ticket, 0)

	generatedId := uuid.New()
	ticket := models.Ticket{
		Id:             generatedId,
		FlightInfo:     models.Flight{},
		Owner:          models.Customer{},
		TicketCost:     150.75,
		ClassOfService: "Economy",
	}
	TicketsCart = append(TicketsCart, ticket)

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

	if len(TicketsCart) != 0 {
		t.Errorf("Expected ticket to be deleted, got %v", TicketsCart)
	}
}

func TestClearHandler(t *testing.T) {
	ticket1 := models.Ticket{
		Id:             uuid.New(),
		FlightInfo:     models.Flight{},
		Owner:          models.Customer{},
		TicketCost:     150.75,
		ClassOfService: "Economy",
	}
	TicketsCart = append(TicketsCart, ticket1)

	ticket2 := models.Ticket{
		Id:             uuid.New(),
		FlightInfo:     models.Flight{},
		Owner:          models.Customer{},
		TicketCost:     200.00,
		ClassOfService: "Business",
	}
	TicketsCart = append(TicketsCart, ticket2)

	req, err := http.NewRequest("DELETE", "/cart/clear", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ClearHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ClearHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedTickets []models.Ticket
	if err := json.Unmarshal(rr.Body.Bytes(), &returnedTickets); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(returnedTickets) != 2 {
		t.Errorf("Expected 2 tickets to be returned, got %v", len(returnedTickets))
	}

	if len(TicketsCart) != 0 {
		t.Errorf("Expected cart to be empty, got %v", len(TicketsCart))
	}
}
