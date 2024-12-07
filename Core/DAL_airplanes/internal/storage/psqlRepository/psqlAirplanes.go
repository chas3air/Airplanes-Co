package psql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/chas3air/Airplanes-Co/Core/DAL_airplanes/internal/config"
	"github.com/chas3air/Airplanes-Co/Core/DAL_airplanes/internal/models"
)

type PsqlAirplanesStorage struct {
	PsqlStorage
}

func MustNewPsqlTicketsStorage() PsqlAirplanesStorage {
	return PsqlAirplanesStorage{
		PsqlStorage: PsqlStorage{},
	}
}

func (s PsqlAirplanesStorage) GetAll(ctx context.Context) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlAirplanes.GetAll"

	log.Println("Using table name:", os.Getenv("PSQL_TABLE_NAME"))

	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	// Use the DB from PsqlStorage instead of opening a new connection
	rows, err := db.QueryContext(ctx, `
		SELECT id, name FROM `+os.Getenv("PSQL_TABLE_NAME")+`;
	`)
	if err != nil {
		log.Println("Error querying airplanes:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}
	defer rows.Close()

	var airplanes []models.Airplane

	for rows.Next() {
		var airplane models.Airplane
		if err := rows.Scan(&airplane.Id, &airplane.Name); err != nil {
			log.Println("Error scanning row:", err.Error())
			continue
		}
		airplanes = append(airplanes, airplane)
	}

	if err := rows.Err(); err != nil {
		log.Println("Row error:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Retrieved %d airplanes\n", len(airplanes))
	return airplanes, nil
}

func (s PsqlAirplanesStorage) Insert(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlAirplanes.Insert"

	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	// Use the injected DB instead of opening a new connection
	airplane := innerObj.(models.Airplane)

	_, err = db.ExecContext(ctx, `
		INSERT INTO `+os.Getenv("PSQL_TABLE_NAME")+`
		(id, name) VALUES ($1, $2)
	`, airplane.Id, airplane.Name)

	if err != nil {
		log.Println("Error inserting airplane:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Inserted airplane with ID=%s\n", airplane.Id)
	return airplane, nil
}

func (s PsqlAirplanesStorage) Update(ctx context.Context, innerObj any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlAirplanes.Update"

	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	airplane := innerObj.(models.Airplane)

	// Use the injected DB instead of opening a new connection
	_, err = db.ExecContext(ctx, `
		UPDATE `+os.Getenv("PSQL_TABLE_NAME")+`
		SET name = $1
		WHERE id = $2;
	`, airplane.Name, airplane.Id)

	if err != nil {
		log.Println("Error updating airplane:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Updated airplane with ID=%s\n", airplane.Id)
	return airplane, nil
}

func (s PsqlAirplanesStorage) Delete(ctx context.Context, id any) (any, error) {
	const op = "DAL.internal.storage.psqlRepository.psqlAirplanes.Delete"

	db, err := sql.Open("postgres", config.ConnStr)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	// Use the injected DB instead of opening a new connection
	row := db.QueryRowContext(ctx, `
		SELECT * FROM `+os.Getenv("PSQL_TABLE_NAME")+`
		WHERE id = $1;
	`, id)

	var airplane models.Airplane
	err = row.Scan(&airplane.Id, &airplane.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No airplane found with ID=%s\n", id)
			return nil, fmt.Errorf("%s: no airplane found with ID=%v", op, id)
		}
		log.Println("Error scanning airplane:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	_, err = db.ExecContext(ctx, `
		DELETE FROM `+os.Getenv("PSQL_TABLE_NAME")+`
		WHERE id = $1;
	`, id)

	if err != nil {
		log.Println("Error deleting airplane:", err.Error())
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	log.Printf("Deleted airplane with ID=%s\n", id)
	return airplane, nil
}
