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

// GetItemHandler handles GET requests to retrieve an item from the cache.
// It expects a URL parameter "key" which specifies the cache key to retrieve.
// If the key is found and not expired, it returns the associated value as JSON.
// If the key is not found or has expired, it returns an appropriate HTTP error.
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

// SetItemHandler handles POST requests to store an item in the cache.
// It expects a JSON request body containing a key and a cacheItem object,
// which includes the value to be stored and its expiration time.
// Upon successful storage, it returns an HTTP 200 status.
// If there is an error reading the body or unmarshaling the JSON, it returns an error.
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

// DeleteItemHandler handles DELETE requests to remove an item from the cache.
// It expects a URL parameter "key" which specifies the cache key to delete.
// If the key is found and successfully deleted, it returns the deleted value as JSON.
// If the key is not found, it returns an HTTP 204 status indicating no content.
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
