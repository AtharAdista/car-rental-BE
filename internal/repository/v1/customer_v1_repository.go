package repository

import (
	"carrental/internal/model/v1"
	"database/sql"
	"fmt"

	errorsManual "carrental/internal/errors"
)

type CustomerV1Repository struct {
	db *sql.DB
}

func NewCustomerV1Repository(db *sql.DB) *CustomerV1Repository {
	return &CustomerV1Repository{db: db}
}

func (r *CustomerV1Repository) CreateCustomer(customer *model.CreateCustomerV1Req) (int, error) {

	var id int

	err := r.db.QueryRow(`
	INSERT INTO customers_v1 (name, nik, phone_number)
	VALUES ($1, $2, $3)
	RETURNING id
	`, customer.Name, customer.NIK, customer.PhoneNumber).Scan(&id)

	if err != nil {
		return -1, fmt.Errorf("failed to create customer: %w", err)
	}

	return id, nil
}

func (r *CustomerV1Repository) FindAllCustomers() ([]model.CustomerV1, error) {

	rows, err := r.db.Query(`
	SELECT id, name, nik, phone_number FROM customers_v1
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	defer rows.Close()

	var customers []model.CustomerV1

	for rows.Next() {
		var customer model.CustomerV1
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.NIK, &customer.PhoneNumber); err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return customers, nil
}

func (r *CustomerV1Repository) FindCustomerById(id int) (*model.CustomerV1, error) {

	customer := &model.CustomerV1{}

	err := r.db.QueryRow(`
	SELECT id, name, nik, phone_number FROM customers_v1 WHERE id=$1
	`, id).Scan(&customer.ID, &customer.Name, &customer.NIK, &customer.PhoneNumber)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsManual.ErrCustomerNotFound
		}
		return nil, err
	}

	return customer, nil
}

func (r *CustomerV1Repository) UpdateCustomerById(id int, req *model.UpdateCustomerV1Req) (*model.CustomerV1, error) {

	query := "UPDATE customers_v1 SET"
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

	if len(args) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	query = query[:len(query)-1]
	query += fmt.Sprintf(" WHERE id = $%d", argIdx)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	var customer *model.CustomerV1

	customer, err = r.FindCustomerById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find customer: %w", err)
	}

	return customer, nil
}

func (r *CustomerV1Repository) DeleteAllCustomers() ([]model.CustomerV1, error) {

	customer, err := r.FindAllCustomers()

	if err != nil {
		return nil, fmt.Errorf("Failed to get customers: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM customers_v1`)

	if err != nil {
		return nil, fmt.Errorf("failed to delete all customers: %w", err)
	}

	return customer, nil
}

func (r *CustomerV1Repository) DeleteCustomerById(id int) (*model.CustomerV1, error) {

	customer, err := r.FindCustomerById(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to get customer: %w", err)
	}

	_, err = r.db.Exec(`DELETE FROM customers_v1 WHERE id=$1`, id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete customer: %w", err)
	}

	return customer, nil
}
