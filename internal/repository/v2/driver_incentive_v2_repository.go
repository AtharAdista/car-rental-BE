package repository

import (
	"carrental/internal/model/v2"
	"database/sql"
	"fmt"

)

type DriverIncentiveV2Repository struct {
	db *sql.DB
}

func NewDriverIncentiveV2Repository(db *sql.DB) *DriverIncentiveV2Repository {
	return &DriverIncentiveV2Repository{db: db}
}


func (r *DriverIncentiveV2Repository) FindAllDriversIncentives() ([]model.DriverIncentiveV2, error) {

	rows, err := r.db.Query(`
	SELECT id, booking_id, incentive FROM drivers_incentives_v2
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find driver incentive: %w", err)
	}

	defer rows.Close()

	var driversIncentives []model.DriverIncentiveV2

	for rows.Next() {
		var driverIncentive model.DriverIncentiveV2
		if err := rows.Scan(&driverIncentive.ID, &driverIncentive.BookingID, &driverIncentive.Incentive); err != nil {
			return nil, fmt.Errorf("failed to scan driver incentive: %w", err)
		}

		driversIncentives = append(driversIncentives, driverIncentive)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return driversIncentives, nil
}

func (r *DriverIncentiveV2Repository) FindDriverIncentiveById(id int) (*model.DriverIncentiveV2, error) {

	driversIncentives := &model.DriverIncentiveV2{}

	err := r.db.QueryRow(`
	SELECT id, booking_id, incentive FROM drivers_incentives_v2 WHERE id=$1
	`, id).Scan(&driversIncentives.ID, &driversIncentives.BookingID, &driversIncentives.Incentive)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("drivers incentive not found")
		}
		return nil, err
	}

	return driversIncentives, nil
}

