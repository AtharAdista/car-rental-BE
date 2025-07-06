package repository

import (
	"carrental/internal/model/v2"
	"database/sql"
	"fmt"

	errorsManual "carrental/internal/errors"
)

type BookingTypeV2Repository struct {
	db *sql.DB
}

func NewBookingTypeV2Repository(db *sql.DB) *BookingTypeV2Repository {
	return &BookingTypeV2Repository{db: db}
}

func (r *BookingTypeV2Repository) CreateBookingType(bookingType *model.CreateBookingTypeV2Req) (int, error) {

	var id int

	err := r.db.QueryRow(`
	INSERT INTO booking_types_v2 (booking_type, description)
	VALUES ($1, $2)
	RETURNING id
	`, bookingType.BookingType, bookingType.Description).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create booking type: %w", err)
	}

	return id, nil
}

func (r *BookingTypeV2Repository) FindAllBookingType() ([]model.BookingTypeV2, error) {

	rows, err := r.db.Query(`
	SELECT id, booking_type, description FROM booking_types_v2
 	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find booking types: %w", err)
	}

	defer rows.Close()

	var bookingTypes []model.BookingTypeV2

	for rows.Next() {
		var bookingType model.BookingTypeV2

		if err := rows.Scan(&bookingType.ID, &bookingType.BookingType, &bookingType.Description); err != nil {
			return nil, fmt.Errorf("failed to scan booking type: %w", err)
		}

		bookingTypes = append(bookingTypes, bookingType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return bookingTypes, nil
}

func (r *BookingTypeV2Repository) FindBookingTypeById(id int) (*model.BookingTypeV2, error) {

	bookingType := &model.BookingTypeV2{}

	err := r.db.QueryRow(`
	SELECT id, booking_type, description FROM booking_types_v2 WHERE id=$1
	`, id).Scan(&bookingType.ID, &bookingType.BookingType, &bookingType.Description)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsManual.ErrBookingTypeNotFound
		}
		return nil, err
	}

	return bookingType, nil
}

func (r *BookingTypeV2Repository) UpdateBookingTypeById(id int, req *model.UpdateBookingTypeV2Req) (*model.BookingTypeV2, error) {

	query := `UPDATE booking_types_v2 SET`
	args := []interface{}{}
	argIdx := 1

	if req.BookingType != nil {
		query += fmt.Sprintf(" booking_type=$%d,", argIdx)
		args = append(args, *req.BookingType)
		argIdx++
	}

	if req.Description != nil {
		query += fmt.Sprintf(" description=$%d,", argIdx)
		args = append(args, *req.Description)
		argIdx++
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	var bookingType *model.BookingTypeV2

	bookingType, err = r.FindBookingTypeById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find booking type: %w", err)
	}

	return bookingType, nil
}

func (r *BookingTypeV2Repository) DeleteAllBookingTypes() ([]model.BookingTypeV2, error) {

	bookingTypes, err := r.FindAllBookingType()

	if err != nil {
		return nil, fmt.Errorf("Failed to get booking types: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM booking_types_v2`)

	if err != nil {
		return nil, fmt.Errorf("failed to delete all booking types: %w", err)
	}

	return bookingTypes, nil
}

func (r *BookingTypeV2Repository) DeleteBookingTypeById(id int) (*model.BookingTypeV2, error) {

	bookingType, err := r.FindBookingTypeById(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to get booking type: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM booking_types_v2 WHERE id=$1`, id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete booking type: %w", err)
	}

	return bookingType, nil
}
