package service

import (
	"carrental/internal/model/v2"
	"carrental/internal/repository/v2"
	"fmt"
	"time"
)

type BookingV2Service struct {
	bookingV2Repository    *repository.BookingV2Repository
	carsV2Repository       *repository.CarsV2Repository
	driverV2Repository     *repository.DriverV2Repository
	customerV2Repository   *repository.CustomerV2Repository
	membershipV2Repository *repository.MembershipV2Repository
}

func NewBookingV2Service(bookingRepo *repository.BookingV2Repository, carsRepo *repository.CarsV2Repository, driverRepo *repository.DriverV2Repository, customerRepo *repository.CustomerV2Repository, membershipRepo *repository.MembershipV2Repository) *BookingV2Service {
	return &BookingV2Service{bookingV2Repository: bookingRepo, carsV2Repository: carsRepo, driverV2Repository: driverRepo, customerV2Repository: customerRepo, membershipV2Repository: membershipRepo}
}

func (s *BookingV2Service) CreateBooking(req *model.CreateBookingV2Req) (int, error) {

	if req.BookingTypeId == 1 {
		req.DriverID = nil
	}

	if req.DriverID == nil {
		req.BookingTypeId = 1
	}


	car, err := s.carsV2Repository.FindCarById(req.CarsID)
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

	if req.DriverID != nil {
		driver, err := s.driverV2Repository.FindDriverById(*req.DriverID) 

		if err != nil {
			return 0, fmt.Errorf("failed to find driver: %w", err)
		}

		req.TotalDriverCost = float64(duration) * driver.DailyCost
	}
	customer, err := s.customerV2Repository.FindCustomerById(req.CustomerID)

	if err != nil {
		return 0, fmt.Errorf("failed to find customer: %w", err)
	}

	req.TotalCost = float64(duration) * car.DailyRent

	if customer.MembershipId != nil {
		membership, err := s.membershipV2Repository.FindMembershipById(*customer.MembershipId)

		if err != nil {
			return 0, fmt.Errorf("failed to find membership: %w", err)
		}

		req.Discount = req.TotalCost * float64(membership.Discount) / 100.0
	}

	id, err := s.bookingV2Repository.CreateBooking(req)

	if err != nil {
		return 0, fmt.Errorf("failed to update stock car: %w", err)
	}

	if err != nil {
		return 0, fmt.Errorf("failed to create booking: %w", err)
	}

	return id, nil
}

func (s *BookingV2Service) GetAllBookings() ([]model.BookingV2, error) {
	bookings, err := s.bookingV2Repository.FindAllBookings()

	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}

	return bookings, nil
}

func (s *BookingV2Service) GetBookingById(id int) (*model.BookingV2, error) {
	booking, err := s.bookingV2Repository.FindBookingById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get car: %w", err)
	}

	return booking, nil
}

func (s *BookingV2Service) UpdateBookingById(id int, req *model.UpdateBookingV2Req) (*model.BookingV2, error) {

	if req.StartRent != nil || req.EndRent != nil || req.CarsID != nil || req.DriverID.IsSet || req.CustomerID != nil {
		booking, err := s.bookingV2Repository.FindBookingById(id)

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

		if req.BookingTypeId == nil {
			req.BookingTypeId = &booking.BookingTypeId
		}

		if req.CustomerID == nil {
			req.CustomerID = &booking.CustomerID
		}

		if !req.DriverID.IsSet {
			req.DriverID.Value = booking.DriverID
			req.DriverID.IsSet = true
		}

		car, err := s.carsV2Repository.FindCarById(*req.CarsID)

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

		if *req.BookingTypeId == 2 {
			driver, err := s.driverV2Repository.FindDriverById(*req.DriverID.Value)
			if err != nil {
				return nil, fmt.Errorf("driver not found: %w", err)
			}

			totalDriverCost := float64(duration) * driver.DailyCost
			req.TotalDriverCost = &totalDriverCost
		}

		customer, err := s.customerV2Repository.FindCustomerById(*req.CustomerID)

		if err != nil {
			return nil, fmt.Errorf("customer not found: %w", err)
		}

		var discount float64
		if customer.MembershipId != nil {
			membership, err := s.membershipV2Repository.FindMembershipById(*customer.MembershipId)
			if err == nil {
				discount = *req.TotalCost * float64(membership.Discount) / 100.0
			}
		}
		req.Discount = &discount
	}

	updatedBooking, err := s.bookingV2Repository.UpdateBookingById(id, req)

	if err != nil {
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	return updatedBooking, nil
}

func (s *BookingV2Service) DeleteAllBookings() ([]model.BookingV2, error) {

	bookings, err := s.bookingV2Repository.DeleteAllBookings()

	if err != nil {
		return nil, fmt.Errorf("failed to delete books: %w", err)
	}

	return bookings, nil
}

func (s *BookingV2Service) DeleteBookingById(id int) (*model.BookingV2, error) {

	booking, err := s.bookingV2Repository.DeleteBookingById(id)

	if err != nil {
		return nil, fmt.Errorf("failed to delete book: %w", err)
	}

	return booking, nil
}

func (s *BookingV2Service) FinishedStatusBooking(id int) (*model.BookingV2, error) {

	booking, err := s.bookingV2Repository.FinishedStatusBooking(id)

	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	return booking, nil
}
