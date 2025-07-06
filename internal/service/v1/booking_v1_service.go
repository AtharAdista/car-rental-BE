package service

import (
	"carrental/internal/model/v1"
	"carrental/internal/repository/v1"
	"fmt"
	"time"
)

type BookingV1Service struct {
	bookingV1Repository *repository.BookingV1Repository
	carsRepository      *repository.CarsV1Repository
}

func NewBookingV1Service(bookingRepo *repository.BookingV1Repository, carsRepo *repository.CarsV1Repository) *BookingV1Service {
	return &BookingV1Service{bookingV1Repository: bookingRepo, carsRepository: carsRepo}
}

func (s *BookingV1Service) CreateBooking(req *model.CreateBookingV1Req) (int, error) {

	car, err := s.carsRepository.FindCarById(req.CarsID)
	if err != nil {
		return 0, fmt.Errorf("failed to find car: %w", err)
	}

	if car.Stock < 1 {
		return 0, fmt.Errorf("car is full booked")
	}

	start := req.StartRent.ToTime()
	end := req.EndRent.ToTime()

	if end.Before(start) {
		return 0, fmt.Errorf("end date cannot be before start date")
	}

	startDate := req.StartRent.ToTime().Truncate(24 * time.Hour)
	endDate := req.EndRent.ToTime().Truncate(24 * time.Hour)

	duration := int(endDate.Sub(startDate).Hours()/24) + 1
	if duration < 1 {
		duration = 1
	}

	req.TotalCost = float64(duration) * car.DailyRent

	id, err := s.bookingV1Repository.CreateBooking(req)

	if err != nil {
		return 0, fmt.Errorf("failed to update stock car: %w", err)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to create booking: %w", err)
	}

	return id, nil
}

func (s *BookingV1Service) GetAllBookings() ([]model.BookingV1, error) {
	bookings, err := s.bookingV1Repository.FindAllBookings()

	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}

	return bookings, nil
}

func (s *BookingV1Service) GetBookingById(id int) (*model.BookingV1, error) {
	booking, err := s.bookingV1Repository.FindBookingById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	return booking, nil
}

func (s *BookingV1Service) UpdateBookingById(id int, req *model.UpdateBookingV1Req) (*model.BookingV1, error) {

	if req.StartRent != nil || req.EndRent != nil || req.CarsID != nil {
		booking, err := s.bookingV1Repository.FindBookingById(id)

		if err != nil {
			return nil, fmt.Errorf("failed to find booking: %w", err)
		}

		if req.StartRent == nil {
			temp := model.NewDateOnly(booking.StartRent)
			req.StartRent = &temp
		}

		if req.EndRent == nil {
			temp := model.NewDateOnly(booking.EndRent)
			req.EndRent = &temp
		}

		if req.CarsID == nil {
			req.CarsID = &booking.CarsID
		}

		car, err := s.carsRepository.FindCarById(*req.CarsID)

		if err != nil {
			return nil, fmt.Errorf("failed to find car: %w", err)
		}

		if car.Stock < 1 && booking.CarsID != *req.CarsID{
			return nil, fmt.Errorf("car is full booked")
		}

		start := req.StartRent.ToTime()
		end := req.EndRent.ToTime()

		if end.Before(start) {
			return nil, fmt.Errorf("end date cannot be before start date")
		}

		duration := int(end.Sub(start).Hours()/24) + 1

		if duration < 1 {
			duration = 1
		}

		total := float64(duration) * car.DailyRent
		req.TotalCost = &total
	}

	updatedBooking, err := s.bookingV1Repository.UpdateBookingById(id, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	return updatedBooking, nil
}

func (s *BookingV1Service) DeleteAllBookings() ([]model.BookingV1, error) {

	bookings, err := s.bookingV1Repository.DeleteAllBookings()

	if err != nil {
		return nil, fmt.Errorf("failed to delete books: %w", err)
	}

	return bookings, nil
}

func (s *BookingV1Service) DeleteBookingById(id int) (*model.BookingV1, error) {

	booking, err := s.bookingV1Repository.DeleteBookingById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete book: %w", err)
	}

	return booking, nil
}

func (s *BookingV1Service) FinishedStatusBooking(id int) (*model.BookingV1, error) {

	booking, err := s.bookingV1Repository.FinishedStatusBooking(id)

	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return booking, nil
}
