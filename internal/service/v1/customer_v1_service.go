package service

import (
	"carrental/internal/model/v1"
	"carrental/internal/repository/v1"
	"fmt"
)

type CustomerV1Service struct {
	customerV1Repository *repository.CustomerV1Repository
}

func NewCustomerV1Service(repo *repository.CustomerV1Repository) *CustomerV1Service {
	return &CustomerV1Service{customerV1Repository: repo}
}

func (s *CustomerV1Service) CreateCustomer(req *model.CreateCustomerV1Req) (int, error) {

	id, err := s.customerV1Repository.CreateCustomer(req)

	if err != nil {
		return 0, fmt.Errorf("failed to create customer: %w", err)
	}

	return id, nil
}

func (s *CustomerV1Service) GetAllCustomers() ([]model.CustomerV1, error) {
	customers, err := s.customerV1Repository.FindAllCustomers()

	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	return customers, nil
}

func (s *CustomerV1Service) GetCustomerById(id int) (*model.CustomerV1, error) {
	customer, err := s.customerV1Repository.FindCustomerById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerV1Service) UpdateCustomerById(id int, req *model.UpdateCustomerV1Req) (*model.CustomerV1, error) {
	customer, err := s.customerV1Repository.UpdateCustomerById(id, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerV1Service) DeleteAllCustomers() ([]model.CustomerV1, error) {

	customers, err := s.customerV1Repository.DeleteAllCustomers()

	if err != nil {
		return nil, fmt.Errorf("failed to delete customers: %w", err)
	}

	return customers, nil
}

func (s *CustomerV1Service) DeleteCustomerById(id int) (*model.CustomerV1, error) {

	customer, err := s.customerV1Repository.DeleteCustomerById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete customer: %w", err)
	}

	return customer, nil
}
