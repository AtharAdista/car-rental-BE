package repository

import (
	"carrental/internal/model/v1"
	"database/sql"
	"fmt"
	"time"

	errorsManual "carrental/internal/errors"
)

type BookingV1Repository struct {
	db *sql.DB
}

func NewBookingV1Repository(db *sql.DB) *BookingV1Repository {
	return &BookingV1Repository{db: db}
}

func (r *BookingV1Repository) CreateBooking(book *model.CreateBookingV1Req) (int, error) {

	tx, err := r.db.Begin()

	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	var id int

	err = tx.QueryRow(`
	INSERT INTO bookings_v1 (customer_id, cars_id, start_rent, end_rent, total_cost)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`, book.CustomerID, book.CarsID, book.StartRent.ToTime(), book.EndRent.ToTime(), book.TotalCost).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to create book: %w", err)
	}

	_, err = tx.Exec(`UPDATE cars_v1 SET stock = stock - 1 WHERE id = $1 AND stock > 0`, book.CarsID)
	if err != nil {
		return 0, fmt.Errorf("failed to update car stock: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (r *BookingV1Repository) FindAllBookings() ([]model.BookingV1, error) {

	rows, err := r.db.Query(`
	SELECT id, customer_id, cars_id, start_rent, end_rent, total_cost, finished FROM bookings_v1
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find book: %w", err)
	}

	defer rows.Close()

	var bookings []model.BookingV1

	for rows.Next() {
		var booking model.BookingV1
		if err := rows.Scan(&booking.ID, &booking.CustomerID, &booking.CarsID, &booking.StartRent, &booking.EndRent, &booking.TotalCost, &booking.Finished); err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}

		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return bookings, nil
}

func (r *BookingV1Repository) FindBookingById(id int) (*model.BookingV1, error) {

	booking := &model.BookingV1{}

	err := r.db.QueryRow(`
	SELECT id, customer_id, cars_id, start_rent, end_rent, total_cost, finished FROM bookings_v1 WHERE id=$1
	`, id).Scan(&booking.ID, &booking.CustomerID, &booking.CarsID, &booking.StartRent, &booking.EndRent, &booking.TotalCost, &booking.Finished)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsManual.ErrBookNotFound
		}

		return nil, err
	}

	return booking, nil
}

func (r *BookingV1Repository) UpdateBookingById(id int, req *model.UpdateBookingV1Req) (*model.BookingV1, error) {

	tx, err := r.db.Begin()

	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	oldBooking, err := r.FindBookingById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing booking: %w", err)
	}

	query := "UPDATE bookings_v1 SET"
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
			_, err = tx.Exec(`UPDATE cars_v1 SET stock = stock + 1 WHERE id = $1`, oldBooking.CarsID)
			if err != nil {
				return nil, fmt.Errorf("failed to update old car stock: %w", err)
			}

			_, err = tx.Exec(`UPDATE cars_v1 SET stock = stock - 1 WHERE id = $1 AND stock > 0`, *req.CarsID)
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

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	var booking *model.BookingV1

	booking, err = r.FindBookingById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find booking: %w", err)
	}

	

	return booking, nil
}

func (r *BookingV1Repository) DeleteAllBookings() ([]model.BookingV1, error) {
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
		_, err := tx.Exec(`UPDATE cars_v1 SET stock = stock + 1 WHERE id = $1`, booking.CarsID)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock for car %d: %w", booking.CarsID, err)
		}
	}

	_, err = tx.Exec(`DELETE FROM bookings_v1`)
	if err != nil {
		return nil, fmt.Errorf("failed to delete bookings: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return bookings, nil
}

func (r *BookingV1Repository) DeleteBookingById(id int) (*model.BookingV1, error) {

	booking, err := r.FindBookingById(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to get book: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE cars_v1 SET stock = stock + 1 WHERE id = $1`, booking.CarsID)

	if err != nil {
		return nil, fmt.Errorf("failed to yupdate car stock: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM bookings_v1 WHERE id=$1`, id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete book: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return booking, nil
}

func (r *BookingV1Repository) FinishedStatusBooking(id int) (*model.BookingV1, error) {

	booking, err := r.FindBookingById(id)

	if booking.Finished {
		return nil, fmt.Errorf("Booking is Finished, cannot change the status", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE bookings_v1 SET finished = $1 WHERE id = $2`, true, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update booking status: %w", err)
	}

	_, err = tx.Exec(`UPDATE cars_v1 SET stock = stock + 1 WHERE id = $1`, booking.CarsID)
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
