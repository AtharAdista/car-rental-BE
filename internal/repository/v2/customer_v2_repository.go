package repository

import (
	"carrental/internal/model/v2"
	"database/sql"
	"fmt"

	errorsManual "carrental/internal/errors"
)

type CustomerV2Repository struct {
	db *sql.DB
}

func NewCustomerV2Repository(db *sql.DB) *CustomerV2Repository {
	return &CustomerV2Repository{db: db}
}

func (r *CustomerV2Repository) CreateCustomer(customer *model.CreateCustomerV2Req) (int, error) {

	var id int
	

	err := r.db.QueryRow(`
	INSERT INTO customers_v2 (name, nik, phone_number, membership_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`, customer.Name, customer.NIK, customer.PhoneNumber, customer.MembershipId).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to create customer: %w", err)
	}

	return id, nil
}

func (r *CustomerV2Repository) FindAllCustomers() ([]model.CustomerV2, error) {

	rows, err := r.db.Query(`
	SELECT id, name, nik, phone_number, membership_id FROM customers_v2
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	defer rows.Close()

	var customers []model.CustomerV2

	for rows.Next() {
		var customer model.CustomerV2
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.NIK, &customer.PhoneNumber, &customer.MembershipId); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return customers, nil
}

func (r *CustomerV2Repository) FindCustomerById(id int) (*model.CustomerV2, error) {

	customer := &model.CustomerV2{}

	err := r.db.QueryRow(`
	SELECT id, name, nik, phone_number, membership_id FROM customers_v2 WHERE id=$1
	`, id).Scan(&customer.ID, &customer.Name, &customer.NIK, &customer.PhoneNumber, &customer.MembershipId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsManual.ErrCustomerNotFound
		}
		return nil, err
	}

	return customer, nil
}

func (r *CustomerV2Repository) UpdateCustomerById(id int, req *model.UpdateCustomerV2Req) (*model.CustomerV2, error) {

	tx, err := r.db.Begin()

	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	var customer *model.CustomerV2

	customer, err = r.FindCustomerById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find customer before update: %w", err)
	}

	query := "UPDATE customers_v2 SET"
	args := []interface{}{}
	argIdx := 1
	hasUpdate := false

	if req.Name != nil {
		query += fmt.Sprintf(" name=$%d,", argIdx)
		args = append(args, *req.Name)
		argIdx++
		hasUpdate = true
	}

	if req.NIK != nil {
		query += fmt.Sprintf(" nik=$%d,", argIdx)
		args = append(args, *req.NIK)
		argIdx++
		hasUpdate = true
	}

	if req.PhoneNumber != nil {
		query += fmt.Sprintf(" phone_number=$%d,", argIdx)
		args = append(args, *req.PhoneNumber)
		argIdx++
		hasUpdate = true
	}

	fmt.Print(req.MembershipId)


	if req.MembershipId.IsSet {
		if req.MembershipId.Value != nil {
			var discountPercentage float64
			membershipID := *req.MembershipId.Value

			err = tx.QueryRow(`
				SELECT discount FROM memberships_v2 WHERE id = $1
			`, membershipID).Scan(&discountPercentage)

			if err != nil {
				return nil, fmt.Errorf("membership not found: %w", err)
			}

			_, err = tx.Exec(`
				UPDATE bookings_v2
				SET discount = total_cost * $1 / 100.0
				WHERE customer_id = $2
			`, discountPercentage, customer.ID)

			if err != nil {
				return nil, fmt.Errorf("failed to update all booking discounts: %w", err)
			}

			query += fmt.Sprintf(" membership_id=$%d,", argIdx)
			args = append(args, membershipID)
			argIdx++
			hasUpdate = true
		} else {

			_, err = tx.Exec(`
				UPDATE bookings_v2
				SET discount = 0
				WHERE customer_id = $1
			`, customer.ID)
			query += fmt.Sprintf(" membership_id=NULL,")
			hasUpdate = true
		}
	}

	if !hasUpdate {
		return nil, fmt.Errorf("no fields to update")
	}

	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	_, err = tx.Exec(query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	customer, err = r.FindCustomerById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	return customer, nil
}

func (r *CustomerV2Repository) DeleteAllCustomers() ([]model.CustomerV2, error) {

	customer, err := r.FindAllCustomers()

	if err != nil {
		return nil, fmt.Errorf("Failed to get customers: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM customers_v2`)

	if err != nil {
		return nil, fmt.Errorf("failed to delete all customers: %w", err)
	}

	return customer, nil
}

func (r *CustomerV2Repository) DeleteCustomerById(id int) (*model.CustomerV2, error) {

	customer, err := r.FindCustomerById(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to get customer: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM customers_v2 WHERE id=$1`, id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete customer: %w", err)
	}

	return customer, nil
}
