package model

type CarV1 struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Stock     int     `json:"stock"`
	DailyRent float64 `json:"daily_rent"`
}

type CreateCarV1Req struct {
	Name      string  `json:"name"`
	Stock     int     `json:"stock"`
	DailyRent float64 `json:"daily_rent"`
}

type UpdateCarV1Req struct {
	Name      *string  `json:"name"`
	Stock     *int     `json:"stock"`
	DailyRent *float64 `json:"daily_rent"`
}

type UpdateStockCarV1Req struct {
	Stock int `json:"stock"`
}
