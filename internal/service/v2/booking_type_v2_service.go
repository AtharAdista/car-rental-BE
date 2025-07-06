package service

import (
	"carrental/internal/model/v2"
	"carrental/internal/repository/v2"
	"fmt"
)

type BookingTypeV2Service struct {
	bookingTypeV2Repository *repository.BookingTypeV2Repository
}

func NewBookingTypeV2Service(repo *repository.BookingTypeV2Repository) *BookingTypeV2Service {
	return &BookingTypeV2Service{bookingTypeV2Repository: repo}
}

func (s *BookingTypeV2Service) CreateBookingType(req *model.CreateBookingTypeV2Req) (int, error) {

	id, err := s.bookingTypeV2Repository.CreateBookingType(req)

	if err != nil {
		return 0, fmt.Errorf("failed to create booking type: %w", err)
	}

	return id, nil
}

func (s *BookingTypeV2Service) GetAllBookingTypes() ([]model.BookingTypeV2, error) {

	bookingTypes, err := s.bookingTypeV2Repository.FindAllBookingType()

	if err != nil {
		return nil, fmt.Errorf("failed to get booking types: %w", err)
	}

	return bookingTypes, nil
}

func (s *BookingTypeV2Service) GetBookingTypeById(id int) (*model.BookingTypeV2, error) {

	bookingType, err := s.bookingTypeV2Repository.FindBookingTypeById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get booking type: %w", err)
	}

	return bookingType, nil
}

func (s *BookingTypeV2Service) UpdateCarById(id int, req *model.UpdateBookingTypeV2Req) (*model.BookingTypeV2, error) {
	bookingType, err := s.bookingTypeV2Repository.UpdateBookingTypeById(id, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update booking type: %w", err)
	}

	return bookingType, nil
}

func (s *BookingTypeV2Service) DeleteAllBookingTypes() ([]model.BookingTypeV2, error) {

	bookingTypes, err := s.bookingTypeV2Repository.DeleteAllBookingTypes()

	if err != nil {
		return nil, fmt.Errorf("failed to delete booking types: %w", err)
	}

	return bookingTypes, nil
}

func (s *BookingTypeV2Service) DeleteBookingTypeById(id int) (*model.BookingTypeV2, error) {

	bookingType, err := s.bookingTypeV2Repository.DeleteBookingTypeById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete booking type: %w", err)
	}

	return bookingType, nil
}
