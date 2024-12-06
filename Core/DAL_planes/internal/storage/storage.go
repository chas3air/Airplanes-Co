package storage

import (
	"fmt"

	"github.com/chas3air/Airplanes-Co/Core/DAL_planes/internal/storage/interfaces"
)

func MustGetInstanceOfPlanessStorage(query string) interfaces.IPlanesRepository {
	const op = "DAL.internal.storage.mustGetInstanceOfTicketsStorage"
	switch query {
	//case "mongo":
	case "psql":
		return psql.MustNewPsqlPlanesStorage()
	default:
		panic(fmt.Errorf("%s: %s", op, "undefined query string"))
	}
}
