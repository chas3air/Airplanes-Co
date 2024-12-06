package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chas3air/Airplanes-Co/Core/DAL_planes/internal/config"
	"github.com/chas3air/Airplanes-Co/Core/DAL_planes/internal/models"
)

type PsqlPlanesStorage struct {
	PsqlStorage
}

func MustNewPsqlPlanesStorage() PsqlPlanesStorage {
	return PsqlPlanesStorage{
		PsqlStorage: PsqlStorage{},
	}
}

func (s PsqlPlanesStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlPlanes.GetAll"

	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, `
		SELECT * FROM`+os.Getenv("PSQL_TABLE_NAME")+`;
	`)
	if err != nil {
		log.Println("Error querying planes:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	var plane models.Plane

	planes := make([]models.Plane, 0, 10)

	for rows.Next() {
		err := rows.Scan(&plane.Id, &plane.Name, &plane.Capacity)
		if err != nil {
			log.Println("Error scanning row:", err.Error())
			continue
		}
		planes = append(planes, plane)
	}

	if err := rows.Err(); err != nil {
		log.Println("Row error:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrieved %d tickets\n", len(planes))
	return planes, nil
}
