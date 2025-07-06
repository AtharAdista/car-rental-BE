package model

type CarV2 struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Stock     int     `json:"stock"`
	DailyRent float64 `json:"daily_rent"`
}

type CreateCarV2Req struct {
	Name      string  `json:"name"`
	Stock     int     `json:"stock"`
	DailyRent float64 `json:"daily_rent"`
}

type UpdateCarV2Req struct {
	Name      *string  `json:"name"`
	Stock     *int     `json:"stock"`
	DailyRent *float64 `json:"daily_rent"`
}
