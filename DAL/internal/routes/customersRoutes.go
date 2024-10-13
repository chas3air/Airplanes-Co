package routes

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/chas3air/Airplanes-Co/DAL/internal/models"
	"github.com/chas3air/Airplanes-Co/DAL/internal/storage"
	"github.com/gorilla/mux"
)

var DB = storage.MustGetInstanceOfCustomersStorage("psql")

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all customers")
	entities, err := DB.GetAll(context.Background())
	if err != nil {
		log.Printf("Error retrieving customers: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	customers, ok := entities.([]models.Customer)
	if !ok {
		log.Println("Invalid data format")
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(customers)
	if err != nil {
		log.Printf("Cannot marshal object: %s\n", err.Error())
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Succesfully fatched customers.")
}

func GetCustomerById(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching customers by id")

	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Printf("Bad request: invalid ID %s", id_s)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	entity, err := DB.GetById(context.Background(), id)
	if err != nil {
		log.Printf("Error retrieving customer by id: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	customer, ok := entity.(models.Customer)
	if !ok {
		log.Println("Invalid data format")
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(customer)
	if err != nil {
		log.Printf("Cannot marshal object: %s\n", err.Error())
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Succesfully fatched customer by id.")
}

func InsertCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Customer insertion")
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		log.Printf("Cannot unmarshal request body to customer: %s\n", err.Error())
		http.Error(w, "Cannot unmarshal object", http.StatusInternalServerError)
		return
	}

	obj, err := DB.Insert(context.Background(), customer)
	if err != nil {
		log.Printf("Error inserting customer with id: %d\n", customer.Id)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Cannot marshal inserted customer: %s\n", err.Error())
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Succesfully inserted customer with id: %d\n", customer.Id)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		log.Printf("Cannot unmarshal request body to customer: %s", err.Error())
		http.Error(w, "Cannot unmarshal object", http.StatusInternalServerError)
		return
	}

	obj, err := DB.Update(context.Background(), customer)
	if err != nil {
		log.Printf("Error updating customer with ID %d: %s", customer.Id, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Cannot marshal updated customer: %s", err.Error())
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully updated customer with ID: %d", customer.Id)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		log.Printf("Cannot unmarshal request body to customer: %s", err.Error())
		http.Error(w, "Cannot unmarshal object", http.StatusInternalServerError)
		return
	}

	obj, err := DB.Delete(context.Background(), customer.Id)
	if err != nil {
		log.Printf("Error deleting customer with ID %d: %s", customer.Id, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Cannot marshal deleted customer: %s", err.Error())
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully deleted customer with ID: %d", customer.Id)
}
