package repository

import (
	"carrental/internal/model/v2"
	"database/sql"
	"fmt"
	"time"

	errorsManual "carrental/internal/errors"
)

type BookingV2Repository struct {
	db *sql.DB
}

func NewBookingV2Repository(db *sql.DB) *BookingV2Repository {
	return &BookingV2Repository{db: db}
}

func (r *BookingV2Repository) CreateBooking(book *model.CreateBookingV2Req) (int, error) {

	tx, err := r.db.Begin()

	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	var id int

	fmt.Println(book.DriverID)

	err = tx.QueryRow(`
	INSERT INTO bookings_v2 (customer_id, cars_id, booking_type_id, driver_id, start_rent, end_rent, total_cost, discount, total_driver_cost)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`, book.CustomerID, book.CarsID, book.BookingTypeId, book.DriverID, book.StartRent.ToTime(), book.EndRent.ToTime(), book.TotalCost, book.Discount, book.TotalDriverCost).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to create book: %w", err)
	}

	_, err = tx.Exec(`UPDATE cars_v2 SET stock = stock - 1 WHERE id = $1 AND stock > 0`, book.CarsID)
	if err != nil {
		return 0, fmt.Errorf("failed to update car stock: %w", err)
	}

	var dailyRent float64
	err = tx.QueryRow(`SELECT daily_rent FROM cars_v2 WHERE id = $1`, book.CarsID).Scan(&dailyRent)
	if err != nil {
		return 0, fmt.Errorf("failed to get car daily rent: %w", err)
	}

	startDate := book.StartRent.ToTime().Truncate(24 * time.Hour)
	endDate := book.EndRent.ToTime().Truncate(24 * time.Hour)

	duration := int(endDate.Sub(startDate).Hours()/24) + 1
	if duration < 1 {
		duration = 1
	}

	incentive := float64(duration) * dailyRent * 0.05

	_, err = tx.Exec(`INSERT INTO drivers_incentives_v2 (booking_id, incentive)
	VALUES ($1, $2)`, id, incentive)
	if err != nil {
		return 0, fmt.Errorf("failed to update car stock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (r *BookingV2Repository) FindAllBookings() ([]model.BookingV2, error) {

	rows, err := r.db.Query(`
	SELECT id, customer_id, cars_id, booking_type_id, driver_id, start_rent, end_rent, total_cost, finished, discount, total_driver_cost FROM bookings_v2
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find book: %w", err)
	}

	defer rows.Close()

	var bookings []model.BookingV2

	for rows.Next() {
		var booking model.BookingV2
		if err := rows.Scan(&booking.ID, &booking.CustomerID, &booking.CarsID, &booking.BookingTypeId, &booking.DriverID, &booking.StartRent, &booking.EndRent, &booking.TotalCost, &booking.Finished, &booking.Discount, &booking.TotalDriverCost); err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return bookings, nil
}

func (r *BookingV2Repository) FindBookingById(id int) (*model.BookingV2, error) {

	booking := &model.BookingV2{}

	err := r.db.QueryRow(`
	SELECT id, customer_id, cars_id, booking_type_id, driver_id, start_rent, end_rent, total_cost, finished, discount, total_driver_cost FROM bookings_v2 WHERE id=$1
	`, id).Scan(&booking.ID, &booking.CustomerID, &booking.CarsID, &booking.BookingTypeId, &booking.DriverID, &booking.StartRent, &booking.EndRent, &booking.TotalCost, &booking.Finished, &booking.Discount, &booking.TotalDriverCost)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsManual.ErrBookNotFound
		}

		return nil, err
	}

	return booking, nil
}

func (r *BookingV2Repository) UpdateBookingById(id int, req *model.UpdateBookingV2Req) (*model.BookingV2, error) {

	tx, err := r.db.Begin()

	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	oldBooking, err := r.FindBookingById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing booking: %w", err)
	}

	query := "UPDATE bookings_v2 SET"
	args := []interface{}{}
	argIdx := 1

	if req.CustomerID != nil {
		query += fmt.Sprintf(" customer_id=$%d,", argIdx)
		args = append(args, *req.CustomerID)
		argIdx++

	}

	if req.CarsID != nil {
		query += fmt.Sprintf(" cars_id=$%d,", argIdx)
		args = append(args, *req.CarsID)
		argIdx++

		if *req.CarsID != oldBooking.CarsID {
			_, err = tx.Exec(`UPDATE cars_v2 SET stock = stock + 1 WHERE id = $1`, oldBooking.CarsID)
			if err != nil {
				return nil, fmt.Errorf("failed to update old car stock: %w", err)
			}

			_, err = tx.Exec(`UPDATE cars_v2 SET stock = stock - 1 WHERE id = $1 AND stock > 0`, *req.CarsID)
			if err != nil {
				return nil, fmt.Errorf("failed to update new car stock: %w", err)
			}
		}
	}

	if req.StartRent != nil {
		query += fmt.Sprintf(" start_rent=$%d,", argIdx)
		args = append(args, time.Time(*req.StartRent))
		argIdx++
	}

	if req.EndRent != nil {
		query += fmt.Sprintf(" end_rent=$%d,", argIdx)
		args = append(args, time.Time(*req.EndRent))
		argIdx++
	}

	if req.TotalCost != nil {
		query += fmt.Sprintf(" total_cost=$%d,", argIdx)
		args = append(args, *req.TotalCost)
		argIdx++
	}

	if req.Finished != nil {
		query += fmt.Sprintf(" finished=$%d,", argIdx)
		args = append(args, *req.Finished)
		argIdx++
	}

	if req.BookingTypeId != nil {
		query += fmt.Sprintf(" booking_type_id=$%d,", argIdx)
		args = append(args, *req.BookingTypeId)
		argIdx++
	}

	if req.DriverID.IsSet {
		if req.DriverID.Value != nil {
			query += fmt.Sprintf(" driver_id=$%d,", argIdx)
			args = append(args, *req.DriverID.Value)
			argIdx++
		} else {
			query += fmt.Sprintf(" driver_id=NULL,")
		}
	}

	if req.Discount != nil {
		query += fmt.Sprintf(" discount=$%d,", argIdx)
		args = append(args, *req.Discount)
		argIdx++
	}

	if req.TotalDriverCost != nil {
		query += fmt.Sprintf(" total_driver_cost=$%d,", argIdx)
		args = append(args, *req.TotalDriverCost)
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
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	if req.StartRent != nil && req.EndRent != nil && req.CarsID != nil {
		start := req.StartRent.ToTime()
		end := req.EndRent.ToTime()

		duration := int(end.Sub(start).Hours()/24) + 1
		if duration < 1 {
			duration = 1
		}

		var dailyRent float64
		err = tx.QueryRow(`SELECT daily_rent FROM cars_v2 WHERE id = $1`, *req.CarsID).Scan(&dailyRent)
		if err != nil {
			return nil, fmt.Errorf("failed to get car daily rent: %w", err)
		}

		incentive := float64(duration) * dailyRent * 0.05

		_, err = tx.Exec(`
			UPDATE drivers_incentives_v2
			SET incentive = $1
			WHERE booking_id = $2
		`, incentive, id)
		if err != nil {
			return nil, fmt.Errorf("failed to update driver incentive: %w", err)
		}
	}

	var booking *model.BookingV2

	booking, err = r.FindBookingById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find booking: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return booking, nil
}

func (r *BookingV2Repository) DeleteAllBookings() ([]model.BookingV2, error) {
	bookings, err := r.FindAllBookings()
	if err != nil {
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, booking := range bookings {
		_, err := tx.Exec(`UPDATE cars_v2 SET stock = stock + 1 WHERE id = $1`, booking.CarsID)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock for car %d: %w", booking.CarsID, err)
		}
	}

	_, err = tx.Exec(`DELETE FROM bookings_v2`)
	if err != nil {
		return nil, fmt.Errorf("failed to delete bookings: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM drivers_incentives_v2`)

	if err != nil {
		return nil, fmt.Errorf("failed to delete driver incentive : %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return bookings, nil
}

func (r *BookingV2Repository) DeleteBookingById(id int) (*model.BookingV2, error) {

	booking, err := r.FindBookingById(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to get book: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE cars_v2 SET stock = stock + 1 WHERE id = $1`, booking.CarsID)

	if err != nil {
		return nil, fmt.Errorf("failed to update car stock: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM bookings_v2 WHERE id=$1`, id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete book: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM drivers_incentives_v2 WHERE booking_id=$1`, id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete driver incentive : %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return booking, nil
}

func (r *BookingV2Repository) FinishedStatusBooking(id int) (*model.BookingV2, error) {

	booking, err := r.FindBookingById(id)

	if booking.Finished {
		return nil, fmt.Errorf("Booking is Finished, cannot change the status", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE bookings_v2 SET finished = $1 WHERE id = $2`, true, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update booking status: %w", err)
	}

	_, err = tx.Exec(`UPDATE cars_v2 SET stock = stock + 1 WHERE id = $1`, booking.CarsID)
	if err != nil {
		return nil, fmt.Errorf("failed to update car stock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	newBooking, err := r.FindBookingById(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to get book: %w", err)
	}

	return newBooking, nil
}
