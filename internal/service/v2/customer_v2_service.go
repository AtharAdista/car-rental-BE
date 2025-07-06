package service

import (
	"carrental/internal/model/v2"
	"carrental/internal/repository/v2"
	"fmt"
)

type CustomerV2Service struct {
	customerV2Repository *repository.CustomerV2Repository
}

func NewCustomerV2Service(repo *repository.CustomerV2Repository) *CustomerV2Service {
	return &CustomerV2Service{customerV2Repository: repo}
}

func (s *CustomerV2Service) CreateCustomer(req *model.CreateCustomerV2Req) (int, error) {

	id, err := s.customerV2Repository.CreateCustomer(req)

	if err != nil {
		return 0, fmt.Errorf("failed to create customer: %w", err)
	}

	return id, nil
}

func (s *CustomerV2Service) GetAllCustomers() ([]model.CustomerV2, error) {
	customers, err := s.customerV2Repository.FindAllCustomers()

	if err != nil {
		return nil, fmt.Errorf("failed to get customers: %w", err)
	}

	return customers, nil
}

func (s *CustomerV2Service) GetCustomerById(id int) (*model.CustomerV2, error) {
	customer, err := s.customerV2Repository.FindCustomerById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerV2Service) UpdateCustomerById(id int, req *model.UpdateCustomerV2Req) (*model.CustomerV2, error) {
	customer, err := s.customerV2Repository.UpdateCustomerById(id, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	return customer, nil
}

func (s *CustomerV2Service) DeleteAllCustomers() ([]model.CustomerV2, error) {

	customers, err := s.customerV2Repository.DeleteAllCustomers()

	if err != nil {
		return nil, fmt.Errorf("failed to delete customers: %w", err)
	}

	return customers, nil
}

func (s *CustomerV2Service) DeleteCustomerById(id int) (*model.CustomerV2, error) {

	customer, err := s.customerV2Repository.DeleteCustomerById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete customer: %w", err)
	}

	return customer, nil
}
