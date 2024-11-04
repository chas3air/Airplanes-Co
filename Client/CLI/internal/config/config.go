package config

import "os"

const DEFAULT_LIMIT_TIME = 5

const (
	FlightsAdmin   string = "flight_admin"
	CustomersAdmin string = "customers_admin"
	GeneralAdmin   string = "general_admin"
	User           string = "user"
	Guest          string = "guest"
)

var Backend_url = os.Getenv("BACKEND_URL")
