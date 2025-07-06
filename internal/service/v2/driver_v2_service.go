package service

import (
	"carrental/internal/model/v2"
	"carrental/internal/repository/v2"
	"fmt"
)

type DriverV2Service struct {
	driverV2Repository *repository.DriverV2Repository
}

func NewDriverV2Service(repo *repository.DriverV2Repository) *DriverV2Service {
	return &DriverV2Service{driverV2Repository: repo}
}

func (s *DriverV2Service) CreateDriver(req *model.CreateDriverV2Req) (int, error) {

	id, err := s.driverV2Repository.CreateDriver(req)

	if err != nil {
		return 0, fmt.Errorf("failed to create driver: %w", err)
	}

	return id, nil
}

func (s *DriverV2Service) GetAllDrivers() ([]model.DriverV2, error) {
	drivers, err := s.driverV2Repository.FindAllDrivers()

	if err != nil {
		return nil, fmt.Errorf("failed to get drivers: %w", err)
	}

	return drivers, nil
}

func (s *DriverV2Service) GetDriverById(id int) (*model.DriverV2, error) {
	driver, err := s.driverV2Repository.FindDriverById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get driver: %w", err)
	}

	return driver, nil
}

func (s *DriverV2Service) UpdateDriverById(id int, req *model.UpdateDriverV2Req) (*model.DriverV2, error) {
	driver, err := s.driverV2Repository.UpdateDriverById(id, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update Driver: %w", err)
	}

	return driver, nil
}

func (s *DriverV2Service) DeleteAllDrivers() ([]model.DriverV2, error) {

	drivers, err := s.driverV2Repository.DeleteAllDrivers()

	if err != nil {
		return nil, fmt.Errorf("failed to delete Driver: %w", err)
	}

	return drivers, nil
}

func (s *DriverV2Service) DeleteDriverById(id int) (*model.DriverV2, error) {

	driver, err := s.driverV2Repository.DeleteDriverById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete membershi: %w", err)
	}

	return driver, nil

}
