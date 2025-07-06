package model

type CustomerV1 struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	NIK         string `json:"nik"`
	PhoneNumber string `json:"phone_number"`
}

type CreateCustomerV1Req struct {
	Name        string `json:"name"`
	NIK         string `json:"nik"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateCustomerV1Req struct {
	Name        *string `json:"name"`
	NIK         *string `json:"nik"`
	PhoneNumber *string `json:"phone_number"`
}
