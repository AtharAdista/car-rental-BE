package model

type DriverIncentiveV2 struct {
	ID        int     `json:"id"`
	BookingID int     `json:"booking_id"`
	Incentive float64 `json:"Incentive"`
}

type CreateDriverIncentiveV2 struct {
	BookingID int     `json:"booking_id"`
	Incentive float64 `json:"Incentive"`
}

type UpdateDriverIncentiveV2 struct {
	BookingID *int     `json:"booking_id"`
	Incentive *float64 `json:"Incentive"`
}
