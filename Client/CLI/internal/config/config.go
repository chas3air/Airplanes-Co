package config

const DEFAULT_LIMIT_TIME = 5

const (
	FlightsAdmin   string = "flight_admin"
	CustomersAdmin string = "customers_admin"
	GeneralAdmin   string = "general_admin"
	User           string = "user"
	Guest          string = "guest"
)

var Backend_url = "http://backend:12013"

//os.Getenv("BACKEND_URL")
