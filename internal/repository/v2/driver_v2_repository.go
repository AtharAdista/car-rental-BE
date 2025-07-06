package repository

import (
	"carrental/internal/model/v2"
	"database/sql"
	"fmt"
	"strings"
	"time"

	errorsManual "carrental/internal/errors"
)

type DriverV2Repository struct {
	db *sql.DB
}

func NewDriverV2Repository(db *sql.DB) *DriverV2Repository {
	return &DriverV2Repository{db: db}
}

func (r *DriverV2Repository) CreateDriver(driver *model.CreateDriverV2Req) (int, error) {

	var id int

	err := r.db.QueryRow(`
	INSERT INTO drivers_v2 (name, nik, phone_number, daily_cost)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`, driver.Name, driver.NIK, driver.PhoneNumber, driver.DailyCost).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create driver: %w", err)
	}

	return id, nil
}

func (r *DriverV2Repository) FindAllDrivers() ([]model.DriverV2, error) {

	rows, err := r.db.Query(`
	SELECT id, name, nik, phone_number, daily_cost FROM drivers_v2
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find drivers: %w", err)
	}

	defer rows.Close()

	var drivers []model.DriverV2

	for rows.Next() {
		var driver model.DriverV2
		if err := rows.Scan(&driver.ID, &driver.Name, &driver.NIK, &driver.PhoneNumber, &driver.DailyCost); err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}

		drivers = append(drivers, driver)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return drivers, nil
}

func (r *DriverV2Repository) FindDriverById(id int) (*model.DriverV2, error) {

	driver := &model.DriverV2{}

	err := r.db.QueryRow(`
	SELECT id, name, nik, phone_number, daily_cost FROM drivers_v2 WHERE id=$1
	`, id).Scan(&driver.ID, &driver.Name, &driver.NIK, &driver.PhoneNumber, &driver.DailyCost)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsManual.ErrDriverNotFound
		}
		return nil, err
	}

	return driver, nil
}

func (r *DriverV2Repository) UpdateDriverById(id int, req *model.UpdateDriverV2Req) (*model.DriverV2, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := "UPDATE drivers_v2 SET"
	args := []interface{}{}
	argIdx := 1

	if req.Name != nil {
		query += fmt.Sprintf(" name=$%d,", argIdx)
		args = append(args, *req.Name)
		argIdx++
	}
	if req.NIK != nil {
		query += fmt.Sprintf(" nik=$%d,", argIdx)
		args = append(args, *req.NIK)
		argIdx++
	}
	if req.PhoneNumber != nil {
		query += fmt.Sprintf(" phone_number=$%d,", argIdx)
		args = append(args, *req.PhoneNumber)
		argIdx++
	}
	if req.DailyCost != nil {
		query += fmt.Sprintf(" daily_cost=$%d,", argIdx)
		args = append(args, *req.DailyCost)
		argIdx++
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query = strings.TrimRight(query, ",")
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	_, err = tx.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update driver: %w", err)
	}

	if req.DailyCost != nil {
		rows, err := tx.Query(`
			SELECT id, start_rent, end_rent
			FROM bookings_v2
			WHERE driver_id = $1 AND finished = false
		`, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get bookings for driver: %w", err)
		}

		var bookings []struct {
			ID        int
			StartRent time.Time
			EndRent   time.Time
		}
		for rows.Next() {
			var b struct {
				ID        int
				StartRent time.Time
				EndRent   time.Time
			}
			if err := rows.Scan(&b.ID, &b.StartRent, &b.EndRent); err != nil {
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
			newCost := float64(days) * float64(*req.DailyCost)

			_, err := tx.Exec(`
				UPDATE bookings_v2
				SET total_driver_cost = $1
				WHERE id = $2 AND finished = false
			`, newCost, b.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to update total_driver_cost: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	driver, err := r.FindDriverById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find updated driver: %w", err)
	}

	return driver, nil
}

func (r *DriverV2Repository) DeleteAllDrivers() ([]model.DriverV2, error) {
	drivers, err := r.FindAllDrivers()
	if err != nil {
		return nil, fmt.Errorf("failed to get drivers: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE bookings_v2
		SET driver_id = NULL,
			total_driver_cost = 0,
			booking_type_id = 1
		WHERE driver_id IS NOT NULL AND finished = false
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to update bookings: %w", err)
	}

	_, err = tx.Exec(`
		DELETE FROM drivers_incentives_v2
		WHERE booking_id NOT IN (
			SELECT id FROM bookings_v2
			WHERE finished = true AND booking_type_id = 2
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to clean up driver incentives: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM drivers_v2`)
	if err != nil {
		return nil, fmt.Errorf("failed to delete all drivers: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return drivers, nil
}

func (r *DriverV2Repository) DeleteDriverById(id int) (*model.DriverV2, error) {
	driver, err := r.FindDriverById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find driver: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE bookings_v2
		SET driver_id = NULL,
			total_driver_cost = 0,
			booking_type_id = 1
		WHERE driver_id = $1 AND finished = false
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update bookings: %w", err)
	}

	_, err = tx.Exec(`
		DELETE FROM drivers_incentives_v2
		WHERE booking_id IN (
			SELECT id FROM bookings_v2
			WHERE driver_id IS NULL AND finished = false
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to delete related driver incentives: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM drivers_v2 WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete driver: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return driver, nil
}
