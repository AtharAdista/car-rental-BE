package model

type CustomerV2 struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	NIK          string `json:"nik"`
	PhoneNumber  string `json:"phone_number"`
	MembershipId *int    `json:"membership_id"`
}

type CreateCustomerV2Req struct {
	Name         string `json:"name"`
	NIK          string `json:"nik"`
	PhoneNumber  string `json:"phone_number"`
	MembershipId *int   `json:"membership_id"`
}

type UpdateCustomerV2Req struct {
	Name         *string `json:"name"`
	NIK          *string `json:"nik"`
	PhoneNumber  *string `json:"phone_number"`
	MembershipId NullableInt  `json:"membership_id"`
}
