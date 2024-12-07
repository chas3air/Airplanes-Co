package storage

import (
	"fmt"

	"github.com/chas3air/Airplanes-Co/Core/DAL_airplanes/internal/storage/interfaces"
	psql "github.com/chas3air/Airplanes-Co/Core/DAL_airplanes/internal/storage/psqlRepository"
)

func MustGetInstanceOfAirplanesStorage(query string) interfaces.IAirplanesRepository {
	const op = "DAL.internal.storage.mustGetInstanceOfAirplanesStorage"
	switch query {
	//case "mongo":
	case "psql":
		return psql.MustNewPsqlTicketsStorage()
	default:
		panic(fmt.Errorf("%s: %s", op, "undefined query string"))
	}
}
