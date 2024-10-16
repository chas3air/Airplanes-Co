package storage

import (
	"fmt"

	"github.com/chas3air/Airplanes-Co/DAL/internal/storage/intefaces"
	psql "github.com/chas3air/Airplanes-Co/DAL/internal/storage/psqlRepository"
)

var db = psql.InitDB()

func MustGetInstanceOfFlightsStorage(query string) intefaces.IFlightsRepository {
	const op = "DAL.internal.storage.mustGetInstanceOfFlightsStorage"
	switch query {
	//case "mongo":
	case "psql":
		return psql.MustNewPsqlFlightsStorage(db)
	default:
		panic(fmt.Errorf("%s: %s", op, "undefined query string"))
	}
}

func MustGetInstanceOfTicketssStorage(query string) intefaces.ITicketsRepository {
	const op = "DAL.internal.storage.mustGetInstanceOfTicketsStorage"
	switch query {
	//case "mongo":
	case "psql":
		return psql.MustNewPsqlTicketsStorage(db)
	default:
		panic(fmt.Errorf("%s: %s", op, "undefined query string"))
	}
}
