package model

type DriverV2 struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	NIK         string  `json:"nik"`
	PhoneNumber string  `json:"phone_number"`
	DailyCost   float64 `json:"daily_cost"`
}

type CreateDriverV2Req struct {
	Name        string  `json:"name"`
	NIK         string  `json:"nik"`
	PhoneNumber string  `json:"phone_number"`
	DailyCost   float64 `json:"daily_cost"`
}

type UpdateDriverV2Req struct {
	Name        *string  `json:"name"`
	NIK         *string  `json:"nik"`
	PhoneNumber *string  `json:"phone_number"`
	DailyCost   *float64 `json:"daily_cost"`
}
