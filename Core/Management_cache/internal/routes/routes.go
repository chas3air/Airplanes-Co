package routes

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chas3air/Airplanes-Co/Core/Management_cache/internal/service"
	"github.com/gorilla/mux"
)

var limiTime = service.GetLimitTime("LIMIT_RESPONSE_TIME")
var cacheURL = os.Getenv("CACHE_URL")

// GetItemHandler retrieves an item from the cache using the provided key.
// It expects a URL parameter "key" that specifies the cache key to retrieve.
// If the key is found, it returns the associated value as JSON.
// If the key is not found or an error occurs during the request, it returns an appropriate HTTP error.
func GetItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get item process...")

	key, ok := mux.Vars(r)["key"]
	if !ok {
		log.Println("Bad request: cannot get key string")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	client := http.Client{
		Timeout: limiTime,
	}

	resp, err := client.Get(cacheURL + "/" + key)
	if err != nil {
		log.Println(err)
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
	log.Printf("Successfully returned item for key: %s\n", key)
}

// SetItemHandler stores an item in the cache.
// It reads the item from the request body, which should be in JSON format.
// If the request is successful, it returns the response from the cache as JSON.
// If an error occurs during the request, it returns an appropriate HTTP error.
func SetItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Set item handler...")

	client := http.Client{
		Timeout: limiTime,
	}

	resp, err := client.Post(cacheURL, "application/json", r.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully set item")
}

// DeleteItemHandler removes an item from the cache using the specified key.
// It expects a URL parameter "key" that specifies the cache key to delete.
// If the key is found and successfully deleted, it returns the deleted item's value as JSON.
// If the key is not found or an error occurs during the request, it returns an appropriate HTTP error.
func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete item process...")

	key, ok := mux.Vars(r)["key"]
	if !ok {
		log.Println("Bad request: cannot get key string")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	req, err := http.NewRequest(http.MethodDelete, cacheURL+"/"+key, nil)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := http.Client{
		Timeout: limiTime,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	log.Println("Successfully deleted item")
}
