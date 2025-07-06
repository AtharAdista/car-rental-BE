package repository

import (
	"carrental/internal/model/v2"
	"database/sql"
	"fmt"

	errorsManual "carrental/internal/errors"
)

type MembershipV2Repository struct {
	db *sql.DB
}

func NewMembershipV2Repository(db *sql.DB) *MembershipV2Repository {
	return &MembershipV2Repository{db: db}
}

func (r *MembershipV2Repository) CreateMembership(membership *model.CreateMembershipV2Req) (int, error) {

	var id int

	err := r.db.QueryRow(`
	INSERT INTO memberships_v2 (membership_name, discount)
	VALUES ($1, $2)
	RETURNING id
	`, membership.MembershipName, membership.Discount).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create membership: %w", err)
	}

	return id, nil
}

func (r *MembershipV2Repository) FindAllMemberships() ([]model.MembershipV2, error) {

	rows, err := r.db.Query(`
	SELECT id, membership_name, discount FROM memberships_v2
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to find memberships: %w", err)
	}

	defer rows.Close()

	var memberships []model.MembershipV2

	for rows.Next() {
		var membership model.MembershipV2
		if err := rows.Scan(&membership.ID, &membership.MembershipName, &membership.Discount); err != nil {
			return nil, fmt.Errorf("failed to scan membership: %w", err)
		}

		memberships = append(memberships, membership)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return memberships, nil
}

func (r *MembershipV2Repository) FindMembershipById(id int) (*model.MembershipV2, error) {

	membership := &model.MembershipV2{}

	err := r.db.QueryRow(`
	SELECT id, membership_name, discount FROM memberships_v2 WHERE id=$1
	`, id).Scan(&membership.ID, &membership.MembershipName, &membership.Discount)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorsManual.ErrMembershipNotFound
		}
		return nil, err
	}

	return membership, nil
}

func (r *MembershipV2Repository) UpdateMembershipById(id int, req *model.UpdateMembershipV2Req) (*model.MembershipV2, error) {

	tx, err := r.db.Begin()

	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	query := "UPDATE memberships_v2 SET"
	args := []interface{}{}
	argIdx := 1

	if req.MembershipName != nil {
		query += fmt.Sprintf(" membership_name=$%d,", argIdx)
		args = append(args, *req.MembershipName)
		argIdx++
	}

	if req.Discount != nil {
		query += fmt.Sprintf(" discount=$%d,", argIdx)
		args = append(args, *req.Discount)
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
		return nil, fmt.Errorf("failed to update membership: %w", err)
	}

	if req.Discount != nil {
		_, err = tx.Exec(`
		UPDATE bookings_v2
		SET discount = total_cost * $1 / 100.0
		WHERE finished = false
		AND customer_id IN (
			SELECT id FROM customers_v2
			WHERE membership_id = $2
		)
		`, *req.Discount, id)

		if err != nil {
			return nil, fmt.Errorf("failed to update booking discounts: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	membership, err := r.FindMembershipById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find membership: %w", err)
	}

	

	return membership, nil
}

func (r *MembershipV2Repository) DeleteAllMembership() ([]model.MembershipV2, error) {
	memberships, err := r.FindAllMemberships()
	if err != nil {
		return nil, fmt.Errorf("failed to get memberships: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE customers_v2 SET membership_id = NULL
		WHERE membership_id IS NOT NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to nullify customer memberships: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE bookings_v2
		SET discount = 0
		WHERE finished = false
		AND EXISTS (
			SELECT 1
			FROM customers_v2 c
			WHERE c.id = bookings_v2.customer_id
			AND c.membership_id IS NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to reset booking discounts: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM memberships_v2`)
	if err != nil {
		return nil, fmt.Errorf("failed to delete memberships: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return memberships, nil
}

func (r *MembershipV2Repository) DeleteMembershipById(id int) (*model.MembershipV2, error) {

	membership, err := r.FindMembershipById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to find membership: %w", err)
	}

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		UPDATE customers_v2 SET membership_id = NULL
		WHERE membership_id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to nullify customer membership_id: %w", err)
	}

	_, err = tx.Exec(`
		UPDATE bookings_v2
		SET discount = 0
		WHERE finished = false
		AND customer_id IN (
			SELECT id FROM customers_v2
			WHERE membership_id IS NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to reset discounts on bookings: %w", err)
	}

	_, err = tx.Exec(`
		DELETE FROM memberships_v2 WHERE id=$1
	`, id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete membership: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return membership, nil
}
