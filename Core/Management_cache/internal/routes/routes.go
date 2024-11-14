package routes

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/chas3air/Airplanes-Co/Core/Management_cache/internal/service"
	"github.com/gorilla/mux"
)

var limiTime = service.GetLimitTime("LIMIT_RESPONSE_TIME")
var cacheURL = os.Getenv("CACHE_URL")

func GetAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get items process...")

	client := http.Client{
		Timeout: limiTime,
	}

	resp, err := client.Get(cacheURL)
	if err != nil {
		log.Println("Error fetching items:", err)
		http.Error(w, "Error fetching items", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body:", err)
		http.Error(w, "Cannot read response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if resp.StatusCode == http.StatusOK {
		w.WriteHeader(http.StatusOK)
		w.Write(body)
		log.Println("Successfully returned items")
	} else if resp.StatusCode == http.StatusNoContent {
		w.WriteHeader(http.StatusNoContent)
		log.Println("No items found in cache")
	} else {
		http.Error(w, "Unexpected error", resp.StatusCode)
	}
}

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

	escapedKEY := url.PathEscape(key)
	rawURL := cacheURL + "/" + escapedKEY

	resp, err := client.Get(rawURL)
	if err != nil {
		log.Println("Error fetching item:", err)
		http.Error(w, "Error fetching item", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body:", err)
		http.Error(w, "Cannot read response body", http.StatusInternalServerError)
		return
	}

	bodyStr := string(body)
	cleanedData := strings.ReplaceAll(bodyStr, "\\", "")
	log.Println("Management-cache send to backend:", cleanedData)

	w.Header().Set("Content-Type", "application/json")
	switch resp.StatusCode {
	case http.StatusOK:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cleanedData)[1 : len(cleanedData)-1])
		log.Printf("Successfully returned item for key: %s\n", key)
	case http.StatusNoContent:
		w.WriteHeader(http.StatusNoContent)
		log.Printf("Item with key: %s not found\n", key)
	default:
		http.Error(w, "Unexpected error", resp.StatusCode)
	}
}

func SetItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Set item handler...")

	client := http.Client{
		Timeout: limiTime,
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Bad request: cannot read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := client.Post(cacheURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Println("Error sending request:", err)
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Cannot read response body:", err)
		http.Error(w, "Cannot read response body", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	switch resp.StatusCode {
	case http.StatusCreated:
		w.WriteHeader(http.StatusCreated)
		w.Write(body)
		log.Println("Successfully set item")
	case http.StatusConflict:
		http.Error(w, "Item already exists", http.StatusConflict)
	default:
		http.Error(w, "Unexpected error", resp.StatusCode)
	}
}

func DeleteItemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Delete item process...")

	key, ok := mux.Vars(r)["key"]
	if !ok {
		log.Println("Bad request: cannot get key string")
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	escapedKEY := url.PathEscape(key)
	rawURL := cacheURL + "/" + escapedKEY
	req, err := http.NewRequest(http.MethodDelete, rawURL, nil)
	if err != nil {
		log.Println("Cannot create request:", err)
		http.Error(w, "Cannot create request", http.StatusInternalServerError)
		return
	}

	client := http.Client{
		Timeout: limiTime,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending delete request:", err)
		http.Error(w, "Error sending delete request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Cannot read response body:", err)
			http.Error(w, "Cannot read response body", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
		log.Println("Successfully deleted item")
	} else {
		http.Error(w, "Unexpected error", resp.StatusCode)
	}
}
