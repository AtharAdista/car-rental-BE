package model

type MembershipV2 struct {
	ID             int    `json:"id"`
	MembershipName string `json:"membership_name"`
	Discount       int    `json:"discount"`
}

type CreateMembershipV2Req struct {
	MembershipName string `json:"membership_name"`
	Discount       int    `json:"discount"`
}

type UpdateMembershipV2Req struct {
	MembershipName *string `json:"membership_name"`
	Discount       *int    `json:"discount"`
}
