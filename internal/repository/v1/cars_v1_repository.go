package repository

import (
	"carrental/internal/model/v1"
	"database/sql"
	"fmt"
	"time"

	errorsManual "carrental/internal/errors"
)

type CarsV1Repository struct {
	db *sql.DB
}

func NewCarsV1Repository(db *sql.DB) *CarsV1Repository {
	return &CarsV1Repository{db: db}
}

func (r *CarsV1Repository) CreateCar(cars *model.CreateCarV1Req) error {

	_, err := r.db.Exec(`
	INSERT INTO cars_v1 (name, stock, daily_rent)
	VALUES ($1, $2, $3)`,
		cars.Name, cars.Stock, cars.DailyRent)

	if err != nil {
		return fmt.Errorf("failed to create cars: %w", err)
	}

	return nil
}

func (r *CarsV1Repository) FindAllCars() ([]model.CarV1, error) {

	rows, err := r.db.Query(`
	SELECT * FROM cars_v1
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find cars: %w", err)
	}

	defer rows.Close()

	var cars []model.CarV1

	for rows.Next() {
		var car model.CarV1
		if err := rows.Scan(&car.ID, &car.Name, &car.Stock, &car.DailyRent); err != nil {
			return nil, fmt.Errorf("failed to scan car: %w", err)
		}

		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return cars, nil
}

func (r *CarsV1Repository) FindCarById(id int) (*model.CarV1, error) {

	car := &model.CarV1{}

	err := r.db.QueryRow(`
	SELECT id, name, stock, daily_rent FROM cars_v1 WHERE id=$1
	`, id).Scan(&car.ID, &car.Name, &car.Stock, &car.DailyRent)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsManual.ErrCarNotFound
		}
		return nil, err
	}

	return car, nil
}

func (r *CarsV1Repository) UpdateCarById(id int, req *model.UpdateCarV1Req) (*model.CarV1, error) {

	tx, err := r.db.Begin()

	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := "UPDATE cars_v1 SET"
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		query += fmt.Sprintf(" name=$%d,", argIdx)
		args = append(args, *req.Name)
		argIdx++
	}

	if req.Stock != nil {
		query += fmt.Sprintf(" stock = $%d,", argIdx)
		args = append(args, *req.Stock)
		argIdx++
	}

	if req.DailyRent != nil {
		rows, err := tx.Query(`
        SELECT id, start_rent, end_rent
        FROM bookings_v1
        WHERE cars_id = $1 AND finished = false
   		 `, id)

		if err != nil {
			return nil, fmt.Errorf("failed to query active bookings: %w", err)
		}

		defer rows.Close()

		for rows.Next() {
			var bookingID int
			var startRent, endRent time.Time

			if err := rows.Scan(&bookingID, &startRent, &endRent); err != nil {
				return nil, fmt.Errorf("failed to scan booking row: %w", err)
			}

			days := int(endRent.Sub(startRent).Hours()/24) + 1

			if days < 1 {
				days = 1
			}

			newTotalCost := float64(days) * float64(*req.DailyRent)

			_, err := tx.Exec(`
				UPDATE bookings_v1
				SET total_cost = $1
				WHERE id = $2
			`, newTotalCost, bookingID)

			if err != nil {
				return nil, fmt.Errorf("failed to update booking cost: %w", err)
			}
		}

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}

		query += fmt.Sprintf(" daily_rent = $%d,", argIdx)
		args = append(args, *req.DailyRent)
		argIdx++
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	_, err = tx.Exec(query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	var car *model.CarV1

	car, err = r.FindCarById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find car: %w", err)
	}

	return car, nil
}

func (r *CarsV1Repository) DeleteAllCars() ([]model.CarV1, error) {

	car, err := r.FindAllCars()

	if err != nil {
		return nil, fmt.Errorf("Failed to get cars: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM cars_v1`)

	if err != nil {
		return nil, fmt.Errorf("failed to delete all cars: %w", err)
	}

	return car, nil
}

func (r *CarsV1Repository) DeleteCarById(id int) (*model.CarV1, error) {

	car, err := r.FindCarById(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to get car: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM cars_v1 WHERE id=$1`, id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete car: %w", err)
	}

	return car, nil
}
