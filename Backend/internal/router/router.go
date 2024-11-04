package router

import (
	"os"
)

var auth_api_url = os.Getenv("AUTHAPI_URL")
var catalog_api_url = os.Getenv("CATALOG_API_URL")
var management_customers_api_url = os.Getenv("MANAGEMENT_CUSTOMERS_API_URL")
var management_flights_api_url = os.Getenv("MANAGEMENT_FLIGHTS_API_URL")
var management_tickets_api_url = os.Getenv("MANAGEMENT_TICKETS_API_URL")
var purchased_tickets_api_url = os.Getenv("PURCHASED_TICKETS_API_URL")
var order_tickets_api_url = os.Getenv("ORDER_TICKET_API_URL")
var cart_api_url = os.Getenv("CART_API_URL")
