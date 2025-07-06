package service

import (
	"carrental/internal/model/v2"
	"carrental/internal/repository/v2"
	"fmt"
)

type MembershipV2Service struct {
	membershipV2Repository *repository.MembershipV2Repository
}

func NewMembershipV2Service(repo *repository.MembershipV2Repository) *MembershipV2Service {
	return &MembershipV2Service{membershipV2Repository: repo}
}

func (s *MembershipV2Service) CreateMembership(req *model.CreateMembershipV2Req) (int, error) {

	id, err := s.membershipV2Repository.CreateMembership(req)

	if err != nil {
		return 0, fmt.Errorf("failed to create membership: %w", err)
	}

	return id, nil
}

func (s *MembershipV2Service) GetAllMemberships() ([]model.MembershipV2, error) {
	memberships, err := s.membershipV2Repository.FindAllMemberships()

	if err != nil {
		return nil, fmt.Errorf("failed to get memberships: %w", err)
	}

	return memberships, nil
}

func (s *MembershipV2Service) GetMembershipById(id int) (*model.MembershipV2, error) {
	membership, err := s.membershipV2Repository.FindMembershipById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}

	return membership, nil
}

func (s *MembershipV2Service) UpdateMembershipById(id int, req *model.UpdateMembershipV2Req) (*model.MembershipV2, error) {
	membership, err := s.membershipV2Repository.UpdateMembershipById(id, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update membership: %w", err)
	}

	return membership, nil
}

func (s *MembershipV2Service) DeleteAllMembership() ([]model.MembershipV2, error) {

	memberships, err := s.membershipV2Repository.DeleteAllMembership()

	if err != nil {
		return nil, fmt.Errorf("failed to delete membership: %w", err)
	}

	return memberships, nil
}

func (s *MembershipV2Service) DeleteMembershipById(id int) (*model.MembershipV2, error) {

	membership, err := s.membershipV2Repository.DeleteMembershipById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete membershi: %w", err)
	}

	return membership, nil

}
