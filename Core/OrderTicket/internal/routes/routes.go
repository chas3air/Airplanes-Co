package routes

import (
	"net/http"

	"github.com/chas3air/Airplanes-Co/Core/Checkbook/internal/service"
)

var limitTime = service.GetLimitTime("LIMIT_RESPONSE_TIME")

var httpClient = &http.Client{
	Timeout: limitTime,
}

func InsertTicketsToCheckbookHandler(w http.ResponseWriter, r *http.Request) {}
func GetTicketsFromCheckbookHandler(w http.ResponseWriter, r *http.Request)  {}
