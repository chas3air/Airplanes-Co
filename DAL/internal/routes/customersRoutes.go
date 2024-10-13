package routes

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/chas3air/Airplanes-Co/DAL/internal/models"
	"github.com/chas3air/Airplanes-Co/DAL/internal/storage"
	"github.com/gorilla/mux"
)

var DB = storage.MustGetInstanceOfCustomersStorage("psql")

func GetCustomers(w http.ResponseWriter, r *http.Request) {
	entities, err := DB.GetAll(context.Background())
	if err != nil {
		http.Error(w, "Server Internal Error", http.StatusInternalServerError)
		return
	}

	customers, ok := entities.([]models.Customer)
	if !ok {
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(customers)
	if err != nil {
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func GetCustomerById(w http.ResponseWriter, r *http.Request) {
	id_s := mux.Vars(r)["id"]
	id, err := strconv.Atoi(id_s)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	entity, err := DB.GetById(context.Background(), id)
	if err != nil {
		http.Error(w, "Server Internal Error", http.StatusInternalServerError)
		return
	}

	customer, ok := entity.(models.Customer)
	if !ok {
		http.Error(w, "Invalid data format", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(customer)
	if err != nil {
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}

func InsertCustomer(w http.ResponseWriter, r *http.Request) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		http.Error(w, "Cannot unmarshal object", http.StatusInternalServerError)
		return
	}

	obj, err := DB.Insert(context.Background(), customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
}

func UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		http.Error(w, "Cannot unmarshal object", http.StatusInternalServerError)
		return
	}

	obj, err := DB.Update(context.Background(), customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
}

func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	bs, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var customer models.Customer
	err = json.Unmarshal(bs, &customer)
	if err != nil {
		http.Error(w, "Cannot unmarshal object", http.StatusInternalServerError)
		return
	}

	obj, err := DB.Delete(context.Background(), customer.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	out_bs, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, "Cannot marshal object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out_bs)
}
