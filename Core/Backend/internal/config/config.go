package config

import "os"

const DEFAULT_LIMIT_TIME = 5

var Management_flights_api_url = os.Getenv("MANAGEMENT_FLIGHTS_API_URL")
var Management_cache_api_url = os.Getenv("MANAGEMENT_CACHE_API_URL")
var Management_customers_api_url = os.Getenv("MANAGEMENT_CUSTOMERS_API_URL")
var Management_tickets_api_url = os.Getenv("MANAGEMENT_TICKETS_API_URL")
var Catalog_api_url = os.Getenv("CATALOG_API_URL")
var Auth_api_url = os.Getenv("AUTHAPI_URL")

const KEY_FOR_FLIGHTS = "flights"
const KEY_FOR_CUSTOMERS = "customers"
const KEY_FOR_TICKETS = "tickets"
const VALUE_EXPIRATION_TIME = 30

var purchased_tickets_api_url = os.Getenv("PURCHASED_TICKETS_API_URL")
var order_tickets_api_url = os.Getenv("ORDER_TICKET_API_URL")
var cart_api_url = os.Getenv("CART_API_URL")
