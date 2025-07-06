package model

import "time"

type BookingV1 struct {
	ID         int       `json:"id"`
	CustomerID int       `json:"customer_id"`
	CarsID     int       `json:"cars_id"`
	StartRent  time.Time `json:"start_rent"`
	EndRent    time.Time `json:"end_rent"`
	TotalCost  float64   `json:"total_cost"`
	Finished   bool      `json:"finished"`
}

type CreateBookingV1Req struct {
	CustomerID int      `json:"customer_id"`
	CarsID     int      `json:"cars_id"`
	StartRent  DateOnly `json:"start_rent"`
	TotalCost  float64  `json:"total_cost"`
	EndRent    DateOnly `json:"end_rent"`
}

type UpdateBookingV1Req struct {
	CustomerID *int      `json:"customer_id"`
	CarsID     *int      `json:"cars_id"`
	StartRent  *DateOnly `json:"start_rent"`
	EndRent    *DateOnly `json:"end_rent"`
	TotalCost  *float64  `json:"total_cost"`
	Finished   *bool     `json:"finished"`
}
