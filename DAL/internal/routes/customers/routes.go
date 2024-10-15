package customers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/chas3air/Airplanes-Co/DAL/internal/config"
	"github.com/chas3air/Airplanes-Co/DAL/internal/models"
	"github.com/chas3air/Airplanes-Co/DAL/internal/storage"
	"github.com/gorilla/mux"
)

var CustomersDB = storage.MustGetInstanceOfCustomersStorage("psql")

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all customers")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

	entities, err := CustomersDB.GetAll(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error retrieving customers: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	customers, ok := entities.([]models.Customer)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(customers)
	if err != nil {
		log.Printf("Cannot marshal object: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Succesfully fetched customers.")
}

func GetCustomerById(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching customers by id")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Printf("Bad request: invalid ID: %s", id_s)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	entity, err := CustomersDB.GetById(context.Background(), id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error retrieving customer by id: %s\n", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	customer, ok := entity.(models.Customer)
	if !ok {
		log.Println("Invalid data type")
		http.Error(w, "Invalid data type", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(customer)
	if err != nil {
		log.Printf("Cannot marshal object: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
	log.Println("Succesfully fatched customer by id.")
}

func InsertCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Customer insertion")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		log.Println("Cannot unmarshal request body to customer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := CustomersDB.Insert(context.Background(), customer)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error inserting customer with id: %d\n", customer.Id)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
	log.Println("Customer updating")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

	bs, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		log.Printf("Cannot unmarshal request body to customer: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	obj, err := CustomersDB.Update(context.Background(), customer)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error updating customer with ID %d: %s", customer.Id, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Cannot marshal updated customer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully updated customer with ID: %d", customer.Id)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("Customers deleting")

	ctx, cancel := context.WithTimeout(context.Background(), config.PSQL_LIMIT_RESPONSE_TIME*time.Second)
	defer cancel()

	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		log.Println("Bad request: wrong flight id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	obj, err := CustomersDB.Delete(context.Background(), id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			log.Printf("Error deleting customer with ID %d: %s", id, err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		log.Println("Cannot marshal deleted customer")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
	log.Printf("Successfully deleted customer with ID: %d", id)
}
