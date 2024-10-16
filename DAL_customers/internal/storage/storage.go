package storage

import (
	"fmt"
	psql "github.com/chas3air/Airplanes-Co/DAL_customers/internal/storage/psqlRepository"
	"github.com/chas3air/Airplanes-Co/DAL_customers/internal/storage/intefaces"
)

var db = psql.InitDB()

func MustGetInstanceOfCustomersStorage(query string) intefaces.ICustomersRepository {
	const op = "DAL.internal.storage.mustGetInstanceOfCustomersStorage"
	switch query {
	//case "mongo":
	case "psql":
		return psql.MustNewPsqlCustomersStorage(db)
	default:
		panic(fmt.Errorf("%s: %s", op, "undefined query string"))
	}
}
