package repository

import (
	"carrental/internal/model/v2"
	"database/sql"
	"fmt"
	"time"

	errorsManual "carrental/internal/errors"
)

type CarsV2Repository struct {
	db *sql.DB
}

func NewCarsV2Repository(db *sql.DB) *CarsV2Repository {
	return &CarsV2Repository{db: db}
}

func (r *CarsV2Repository) CreateCar(cars *model.CreateCarV2Req) error {

	_, err := r.db.Exec(`
	INSERT INTO cars_v2 (name, stock, daily_rent)
	VALUES ($1, $2, $3)`,
		cars.Name, cars.Stock, cars.DailyRent)

	if err != nil {
		return fmt.Errorf("failed to create cars: %w", err)
	}

	return nil
}

func (r *CarsV2Repository) FindAllCars() ([]model.CarV2, error) {

	rows, err := r.db.Query(`
	SELECT * FROM cars_v2
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find cars: %w", err)
	}

	defer rows.Close()

	var cars []model.CarV2

	for rows.Next() {
		var car model.CarV2
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

func (r *CarsV2Repository) FindCarById(id int) (*model.CarV2, error) {

	car := &model.CarV2{}

	err := r.db.QueryRow(`
	SELECT id, name, stock, daily_rent FROM cars_v2 WHERE id=$1
	`, id).Scan(&car.ID, &car.Name, &car.Stock, &car.DailyRent)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsManual.ErrCarNotFound
		}
		return nil, err
	}

	return car, nil
}

func (r *CarsV2Repository) UpdateCarById(id int, req *model.UpdateCarV2Req) (*model.CarV2, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := "UPDATE cars_v2 SET"
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		query += fmt.Sprintf(" name=$%d,", argIdx)
		args = append(args, *req.Name)
		argIdx++
	}

	if req.Stock != nil {
		query += fmt.Sprintf(" stock=$%d,", argIdx)
		args = append(args, *req.Stock)
		argIdx++
	}

	if req.DailyRent != nil {
		rows, err := tx.Query(`
			SELECT id, start_rent, end_rent, customer_id
			FROM bookings_v2
			WHERE cars_id = $1 AND finished = false
		`, id)
		if err != nil {
			return nil, fmt.Errorf("failed to query active bookings: %w", err)
		}

		var bookings []struct {
			ID         int
			StartRent  time.Time
			EndRent    time.Time
			CustomerID int
		}
		for rows.Next() {
			var b struct {
				ID         int
				StartRent  time.Time
				EndRent    time.Time
				CustomerID int
			}
			if err := rows.Scan(&b.ID, &b.StartRent, &b.EndRent, &b.CustomerID); err != nil {
				rows.Close()
				return nil, fmt.Errorf("failed to scan booking row: %w", err)
			}
			bookings = append(bookings, b)
		}
		rows.Close()

		for _, b := range bookings {
			days := int(b.EndRent.Sub(b.StartRent).Hours()/24) + 1
			if days < 1 {
				days = 1
			}

			newTotalCost := float64(days) * float64(*req.DailyRent)

			_, err := tx.Exec(`
				UPDATE bookings_v2
				SET total_cost = $1
				WHERE id = $2 and finished=false
			`, newTotalCost, b.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to update booking cost: %w", err)
			}

			var discountPercentage float64
			err = tx.QueryRow(`
				SELECT m.discount
				FROM customers_v2 c
				JOIN memberships_v2 m ON c.membership_id = m.id
				WHERE c.id = $1
			`, b.CustomerID).Scan(&discountPercentage)

			if err != nil {
				discountPercentage = 0 
			}

			newDiscount := newTotalCost * (discountPercentage / 100.0)

			_, err = tx.Exec(`
				UPDATE bookings_v2
				SET discount = $1
				WHERE id = $2 and finished=false
			`, newDiscount, b.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to update discount: %w", err)
			}
		}

		query += fmt.Sprintf(" daily_rent=$%d,", argIdx)
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

	car, err := r.FindCarById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find car: %w", err)
	}

	return car, nil
}

func (r *CarsV2Repository) DeleteAllCars() ([]model.CarV2, error) {

	car, err := r.FindAllCars()

	if err != nil {
		return nil, fmt.Errorf("Failed to get cars: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM cars_v2`)

	if err != nil {
		return nil, fmt.Errorf("failed to delete all cars: %w", err)
	}

	return car, nil
}

func (r *CarsV2Repository) DeleteCarById(id int) (*model.CarV2, error) {

	car, err := r.FindCarById(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to get car: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM cars_v2 WHERE id=$1`, id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete car: %w", err)
	}

	return car, nil
}
