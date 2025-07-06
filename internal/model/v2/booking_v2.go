package model

import "time"

type BookingV2 struct {
	ID              int       `json:"id"`
	CustomerID      int       `json:"customer_id"`
	CarsID          int       `json:"cars_id"`
	BookingTypeId   int       `json:"booking_type_id"`
	DriverID        *int       `json:"driver_id"`
	StartRent       time.Time `json:"start_rent"`
	EndRent         time.Time `json:"end_rent"`
	TotalCost       float64   `json:"total_cost"`
	Finished        bool      `json:"finished"`
	Discount        float64   `json:"discount"`
	TotalDriverCost float64   `json:"total_driver_cost"`
}

type CreateBookingV2Req struct {
	CustomerID      int      `json:"customer_id"`
	CarsID          int      `json:"cars_id"`
	BookingTypeId   int      `json:"booking_type_id"`
	DriverID        *int      `json:"driver_id"`
	StartRent       DateOnly `json:"start_rent"`
	EndRent         DateOnly `json:"end_rent"`
	TotalCost       float64  `json:"total_cost"`
	Discount        float64  `json:"discount"`
	TotalDriverCost float64  `json:"total_driver_cost"`
}

type UpdateBookingV2Req struct {
	CustomerID      *int      `json:"customer_id"`
	CarsID          *int      `json:"cars_id"`
	BookingTypeId   *int      `json:"booking_type_id"`
	DriverID        NullableInt   `json:"driver_id"`
	StartRent       *DateOnly `json:"start_rent"`
	EndRent         *DateOnly `json:"end_rent"`
	TotalCost       *float64  `json:"total_cost"`
	Finished        *bool     `json:"finished"`
	Discount        *float64  `json:"discount"`
	TotalDriverCost *float64  `json:"total_driver_cost"`
}
