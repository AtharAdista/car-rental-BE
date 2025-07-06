package service

import (
	"carrental/internal/model/v1"
	"carrental/internal/repository/v1"
	"fmt"
)

type CarsV1Service struct {
	carsV1Repository *repository.CarsV1Repository
}

func NewCarsV1Service(repo *repository.CarsV1Repository) *CarsV1Service {
	return &CarsV1Service{carsV1Repository: repo}
}

func (s *CarsV1Service) CreateCar(req *model.CreateCarV1Req) error {

	if req.Stock < 1 {
		return fmt.Errorf("stock must greater than 0")
	}

	err := s.carsV1Repository.CreateCar(req)

	if err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}

	return nil
}

func (s *CarsV1Service) GetAllCars() ([]model.CarV1, error) {
	cars, err := s.carsV1Repository.FindAllCars()

	if err != nil {
		return nil, fmt.Errorf("failed to get cars: %w", err)
	}

	return cars, nil
}

func (s *CarsV1Service) GetCarById(id int) (*model.CarV1, error) {
	car, err := s.carsV1Repository.FindCarById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	return car, nil
}

func (s *CarsV1Service) UpdateCarById(id int, req *model.UpdateCarV1Req) (*model.CarV1, error) {

	if req.Stock != nil {
		if *req.Stock < 1 {
		return nil, fmt.Errorf("stock must greater than 0")
		}	
	}

	car, err := s.carsV1Repository.UpdateCarById(id, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	return car, nil
}

func (s *CarsV1Service) DeleteAllCars() ([]model.CarV1, error) {

	cars, err := s.carsV1Repository.DeleteAllCars()

	if err != nil {
		return nil, fmt.Errorf("failed to delete cars: %w", err)
	}

	return cars, nil
}

func (s *CarsV1Service) DeleteCarById(id int) (*model.CarV1, error) {

	car, err := s.carsV1Repository.DeleteCarById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete cars: %w", err)
	}

	return car, nil
}
