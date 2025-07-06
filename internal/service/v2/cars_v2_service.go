package service

import (
	"carrental/internal/model/v2"
	"carrental/internal/repository/v2"
	"fmt"
)

type CarsV2Service struct {
	carsV2Repository *repository.CarsV2Repository
}

func NewCarsV2Service(repo *repository.CarsV2Repository) *CarsV2Service {
	return &CarsV2Service{carsV2Repository: repo}
}

func (s *CarsV2Service) CreateCar(req *model.CreateCarV2Req) error {

	if req.Stock < 1 {
		return fmt.Errorf("stock must greater than 0")
	}

	err := s.carsV2Repository.CreateCar(req)

	if err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	return nil
}

func (s *CarsV2Service) GetAllCars() ([]model.CarV2, error) {
	cars, err := s.carsV2Repository.FindAllCars()

	if err != nil {
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}

	return cars, nil
}

func (s *CarsV2Service) GetCarById(id int) (*model.CarV2, error) {
	car, err := s.carsV2Repository.FindCarById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	return car, nil
}

func (s *CarsV2Service) UpdateCarById(id int, req *model.UpdateCarV2Req) (*model.CarV2, error) {

	if req.Stock != nil {
		if *req.Stock < 1 {
		return nil, fmt.Errorf("stock must greater than 0")
		}	
	}

	car, err := s.carsV2Repository.UpdateCarById(id, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	return car, nil
}

func (s *CarsV2Service) DeleteAllCars() ([]model.CarV2, error) {

	cars, err := s.carsV2Repository.DeleteAllCars()

	if err != nil {
		return nil, fmt.Errorf("failed to delete cars: %w", err)
	}

	return cars, nil
}

func (s *CarsV2Service) DeleteCarById(id int) (*model.CarV2, error) {

	car, err := s.carsV2Repository.DeleteCarById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete cars: %w", err)
	}

	return car, nil
}
