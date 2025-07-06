package service

import (
	"carrental/internal/model/v2"
	"carrental/internal/repository/v2"
	"fmt"
)

type DriverIncentiveV2Service struct {
	driverIncentiveV2Repository *repository.DriverIncentiveV2Repository
}

func NewDriverIncentiveV2Service(repo *repository.DriverIncentiveV2Repository) *DriverIncentiveV2Service {
	return &DriverIncentiveV2Service{driverIncentiveV2Repository: repo}
}

func (s *DriverIncentiveV2Service) GetAllDriverIncentives() ([]model.DriverIncentiveV2, error) {
	driverIncentives, err := s.driverIncentiveV2Repository.FindAllDriversIncentives()

	if err != nil {
		return nil, fmt.Errorf("failed to get driverIncentives: %w", err)
	}

	return driverIncentives, nil
}

func (s *DriverIncentiveV2Service) GetDriverIncentiveById(id int) (*model.DriverIncentiveV2, error) {
	driverIncentive, err := s.driverIncentiveV2Repository.FindDriverIncentiveById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get driverIncentive: %w", err)
	}

	return driverIncentive, nil
}