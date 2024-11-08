package router

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chas3air/Airplanes-Co/Core/Cart/internal/models"
	"github.com/google/uuid"
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
		FlightId:       uuid.New(),
		OwnerId:        uuid.New(),
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
	if len(TicketsCart) != 1 || TicketsCart[0] != ticket {
		t.Errorf("Expected ticket not found in cart: got %v", TicketsCart)
	}
}

func TestUpdateTicketHandler(t *testing.T) {
	// Сначала добавим билет, чтобы его обновить
	ticket := models.Ticket{
		Id:             uuid.New(),
		FlightId:       uuid.New(),
		OwnerId:        uuid.New(),
		TicketCost:     150.75,
		ClassOfService: "Economy",
	}
	TicketsCart = append(TicketsCart, ticket)

	// Обновим билет
	updatedTicket := ticket
	updatedTicket.TicketCost = 200.00
	body, _ := json.Marshal(updatedTicket)

	req, err := http.NewRequest("PATCH", "/cart", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UpdateTicketHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("UpdateTicketHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем, что билет был обновлён
	if TicketsCart[0].TicketCost != 200.00 {
		t.Errorf("Expected ticket cost to be updated: got %v", TicketsCart[0].TicketCost)
	}
}

func TestDeleteTicketHandler(t *testing.T) {
	// Инициализация TicketsCart перед каждым тестом
	TicketsCart = make([]models.Ticket, 0)

	// Создаем билет и добавляем его в корзину
	generatedId := uuid.New()
	ticket := models.Ticket{
		Id:             generatedId,
		FlightId:       uuid.New(),
		OwnerId:        uuid.New(),
		TicketCost:     150.75,
		ClassOfService: "Economy",
	}
	TicketsCart = append(TicketsCart, ticket)

	log.Println("Generated ID:", generatedId.String())

	// Создаем DELETE запрос с правильным ID
	req, err := http.NewRequest(http.MethodDelete, "/cart/"+generatedId.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем новый Recorder для записи ответов
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteTicketHandler)

	// Обрабатываем запрос
	handler.ServeHTTP(rr, req)

	// Логируем статус ответа
	log.Println("Response status code:", rr.Code)

	// Проверяем, что статус ответа - 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DeleteTicketHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем, что билет был удален
	if len(TicketsCart) != 0 {
		t.Errorf("Expected ticket to be deleted, got %v", TicketsCart)
	}
}

func TestClearHandler(t *testing.T) {
	// Инициализация TicketsCart перед тестом
	//TicketsCart = make([]models.Ticket, 0)

	// Создаем и добавляем билеты в корзину
	ticket1 := models.Ticket{
		Id:             uuid.New(),
		FlightId:       uuid.New(),
		OwnerId:        uuid.New(),
		TicketCost:     150.75,
		ClassOfService: "Economy",
	}
	TicketsCart = append(TicketsCart, ticket1)

	ticket2 := models.Ticket{
		Id:             uuid.New(),
		FlightId:       uuid.New(),
		OwnerId:        uuid.New(),
		TicketCost:     200.00,
		ClassOfService: "Business",
	}
	TicketsCart = append(TicketsCart, ticket2)

	// Создаем DELETE запрос для очистки корзины
	req, err := http.NewRequest("DELETE", "/cart/clear", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ClearHandler)

	// Обрабатываем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем, что статус ответа - 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ClearHandler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Проверяем, что билеты были возвращены
	var returnedTickets []models.Ticket
	if err := json.Unmarshal(rr.Body.Bytes(), &returnedTickets); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(returnedTickets) != 2 {
		t.Errorf("Expected 2 tickets to be returned, got %v", len(returnedTickets))
	}

	// Проверяем, что корзина пуста
	if len(TicketsCart) != 0 {
		t.Errorf("Expected cart to be empty, got %v", len(TicketsCart))
	}
}
