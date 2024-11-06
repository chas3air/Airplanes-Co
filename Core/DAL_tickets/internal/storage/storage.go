package storage

import (
	"fmt"

	"github.com/chas3air/Airplanes-Co/Core/DAL_tickets/internal/storage/intefaces"
	psql "github.com/chas3air/Airplanes-Co/Core/DAL_tickets/internal/storage/psqlRepository"
)

var db = psql.InitDB()

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
