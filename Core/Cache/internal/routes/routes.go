package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Cache/internal/models"
	"github.com/gorilla/mux"
)

var cache = models.NewCarrotCache()

func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Cache get process initiated...")
	key, ok := mux.Vars(r)["key"]
	if !ok {
		log.Println("Bad request: cannot get key string")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to get item from cache with key: %s\n", key)
	value, ok := cache.Get(key)
	if !ok {
		log.Printf("Key not found in cache: %s\n", key)
		http.Error(w, "Cannot find element in cache", http.StatusNoContent)
		return
	}

	log.Printf("Item found in cache: %s\n", key)
	responseData, err := json.Marshal(value)
	if err != nil {
		log.Println("Error marshaling value to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
	log.Printf("Successfully returned item for key: %s\n", key)
}

func SetItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Cache set process initiated...")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := models.Message{}
	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Println("Cannot unmarshal item:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Setting item in cache with key: %s, expiration: %d\n", message.Key, message.Value.Expiration)
	cache.Set(message.Key, message.Value.Value, message.Value.Expiration)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Printf("Successfully set item in cache with key: %s\n", message.Key)
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Cache delete process initiated...")

	key, ok := mux.Vars(r)["key"]
	if !ok {
		log.Println("Bad request: cannot get key string")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to delete item from cache with key: %s\n", key)
	value := cache.Delete(key)
	if value == nil {
		log.Printf("No item found to delete for key: %s\n", key)
		http.Error(w, "No item found to delete", http.StatusNoContent)
		return
	}

	log.Printf("Item deleted from cache: %s\n", key)
	responseData, err := json.Marshal(value)
	if err != nil {
		log.Println("Error marshaling value to JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	log.Printf("Successfully deleted item for key: %s\n", key)
	w.Write(responseData)
}
