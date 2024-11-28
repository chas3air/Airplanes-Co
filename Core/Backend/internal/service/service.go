package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/config"
	"github.com/chas3air/Airplanes-Co/Core/Backend/internal/models"
)

func GetLimitTime(name string) time.Duration {
	limitTimeEnvS := os.Getenv(name)
	limitTimeEnv, err := strconv.Atoi(limitTimeEnvS)
	if err != nil {
		return time.Duration(config.DEFAULT_LIMIT_TIME)
	}
	return time.Duration(limitTimeEnv) * time.Second
}

func SaveToCache(message models.Message) {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post(config.Management_cache_api_url, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Println("Error posting to cache:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error posting to cache, response code:", resp.StatusCode)
		return
	}

	log.Println("Element successfully seted.")
}
