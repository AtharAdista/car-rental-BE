package model

type BookingTypeV2 struct {
	ID          int    `json:"id"`
	BookingType string `json:"booking_type"`
	Description string `json:"description"`
}

type CreateBookingTypeV2Req struct {
	BookingType string `json:"booking_type"`
	Description string `json:"description"`
}

type UpdateBookingTypeV2Req struct {
	BookingType *string `json:"booking_type"`
	Description *string `json:"description"`
}
