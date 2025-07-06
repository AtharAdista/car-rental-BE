package main

import (
	"carrental/internal/db"
	handlerV1  "carrental/internal/handler/v1"
	handlerV2  "carrental/internal/handler/v2"

	repositoryV1 "carrental/internal/repository/v1"
	repositoryV2 "carrental/internal/repository/v2"

	serviceV1 "carrental/internal/service/v1"
	serviceV2 "carrental/internal/service/v2"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("ENV") != "production" {
		err := godotenv.Load("../../.env")

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	db.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	carsV1Repo := repositoryV1.NewCarsV1Repository(db.GetDB())
	carsV1Service := serviceV1.NewCarsV1Service(carsV1Repo)
	carsV1Handler := handlerV1.NewCarsV1Handler(carsV1Service)

	customerV1Repo := repositoryV1.NewCustomerV1Repository(db.GetDB())
	customerV1Service := serviceV1.NewCustomerV1Service(customerV1Repo)
	customerV1Handler := handlerV1.NewCustomerV1Handler(customerV1Service)

	bookingV1Repo := repositoryV1.NewBookingV1Repository(db.GetDB())
	bookingV1Service := serviceV1.NewBookingV1Service(bookingV1Repo, carsV1Repo)
	bookingV1Handler := handlerV1.NewBookingV1Handler(bookingV1Service)

	carsV2Repo := repositoryV2.NewCarsV2Repository(db.GetDB())
	carsV2Service := serviceV2.NewCarsV2Service(carsV2Repo)
	carsV2Handler := handlerV2.NewCarsV2Handler(carsV2Service)

	driverV2Repo := repositoryV2.NewDriverV2Repository(db.GetDB())
	driverV2Service := serviceV2.NewDriverV2Service(driverV2Repo)
	driverV2Handler := handlerV2.NewDriverV2Handler(driverV2Service)

	membershipV2Repo := repositoryV2.NewMembershipV2Repository(db.GetDB())
	membershipV2Service := serviceV2.NewMembershipV2Service(membershipV2Repo)
	membershipV2Handler := handlerV2.NewMembershipV2Handler(membershipV2Service)

	bookingTypeV2Repo := repositoryV2.NewBookingTypeV2Repository(db.GetDB())
	bookingTypeV2Service := serviceV2.NewBookingTypeV2Service(bookingTypeV2Repo)
	bookingTypeV2Handler := handlerV2.NewBookingTypeV2Handler(bookingTypeV2Service)

	customerV2Repo := repositoryV2.NewCustomerV2Repository(db.GetDB())
	customerV2Service := serviceV2.NewCustomerV2Service(customerV2Repo)
	customerV2Handler := handlerV2.NewCustomerV2Handler(customerV2Service)

	bookingV2Repo := repositoryV2.NewBookingV2Repository(db.GetDB())
	bookingV2Service := serviceV2.NewBookingV2Service(bookingV2Repo, carsV2Repo, driverV2Repo, customerV2Repo, membershipV2Repo)
	bookingV2Handler := handlerV2.NewBookingV2Handler(bookingV2Service)

	driverIncentiveV2Repo := repositoryV2.NewDriverIncentiveV2Repository(db.GetDB())
	driverIncentiveV2Service := serviceV2.NewDriverIncentiveV2Service(driverIncentiveV2Repo)
	driverIncentiveV2Handler := handlerV2.NewDriverIncentiveV2Handler(driverIncentiveV2Service)

	r := gin.Default()

	// V1
	v1 := r.Group("/api/v1")

	// Cars route
	v1.POST("/car", carsV1Handler.CreateCar)
	v1.GET("/cars", carsV1Handler.GetAllCars)
	v1.GET("/car/:id", carsV1Handler.GetCarById)
	v1.PATCH("/car/:id", carsV1Handler.UpdateCarById)
	v1.DELETE("/cars", carsV1Handler.DeleteAllCars)
	v1.DELETE("/car/:id", carsV1Handler.DeleteCarById)

	// Customer route
	v1.POST("/customer", customerV1Handler.CreateCustomer)
	v1.GET("/customers", customerV1Handler.GetAllCustomers)
	v1.GET("/customer/:id", customerV1Handler.GetCustomerById)
	v1.PATCH("/customer/:id", customerV1Handler.UpdateCustomerById)
	v1.DELETE("/customers", customerV1Handler.DeleteAllCustomers)
	v1.DELETE("/customer/:id", customerV1Handler.DeleteCustomerById)

	// Booking route
	v1.POST("/book", bookingV1Handler.CreateBooking)
	v1.GET("/books", bookingV1Handler.GetAllBookings)
	v1.GET("/book/:id", bookingV1Handler.GetBookingById)
	v1.PATCH("/book/:id", bookingV1Handler.UpdateBookingById)
	v1.PATCH("/book/finished/:id", bookingV1Handler.FinishedStatusBooking)
	v1.DELETE("/books", bookingV1Handler.DeleteAllbookings)
	v1.DELETE("/book/:id", bookingV1Handler.DeleteBookingById)

	// V2
	v2 := r.Group("/api/v2")

	// Car route
	v2.POST("/car", carsV2Handler.CreateCar)
	v2.GET("/cars", carsV2Handler.GetAllCars)
	v2.GET("/car/:id", carsV2Handler.GetCarById)
	v2.PATCH("/car/:id", carsV2Handler.UpdateCarById)
	v2.DELETE("/cars", carsV2Handler.DeleteAllCars)
	v2.DELETE("/car/:id", carsV2Handler.DeleteCarById)

	// Driver route
	v2.POST("/driver", driverV2Handler.CreateDriver)
	v2.GET("/drivers", driverV2Handler.GetAllDrivers)
	v2.GET("/driver/:id", driverV2Handler.GetDriverById)
	v2.PATCH("/driver/:id", driverV2Handler.UpdateDriverById)
	v2.DELETE("/drivers", driverV2Handler.DeleteAllDrivers)
	v2.DELETE("/driver/:id", driverV2Handler.DeleteDriverById)

	// Membership route
	v2.POST("/membership", membershipV2Handler.CreateMembership)
	v2.GET("/memberships", membershipV2Handler.GetAllMemberships)
	v2.GET("/membership/:id", membershipV2Handler.GetMembershipById)
	v2.PATCH("/membership/:id", membershipV2Handler.UpdateMembershipById)
	v2.DELETE("/memberships", membershipV2Handler.DeleteAllMemberships)
	v2.DELETE("/membership/:id", membershipV2Handler.DeleteMembershipById)

	// Booking Type route 
	// HATI-HATI KALAU PAKE DELETE KARENA LOGIC KODENYA HANYA CHECK BERDASARKAN
	// ID, yaitu 1 untuk car only dan 2 untuk car and driver
	v2.POST("/booking-type", bookingTypeV2Handler.CreateBookingType)
	v2.GET("/booking-types", bookingTypeV2Handler.GetAllBookingTypes)
	v2.GET("/booking-type/:id", bookingTypeV2Handler.GetBookingTypeById)
	v2.PATCH("/booking-type/:id", bookingTypeV2Handler.UpdateBookingTypeById)
	v2.DELETE("/booking-types", bookingTypeV2Handler.DeleteAllBookingTypes)
	v2.DELETE("/booking-type/:id", bookingTypeV2Handler.DeleteBookingTypeById)

	// Customer route
	v2.POST("/customer", customerV2Handler.CreateCustomer)
	v2.GET("/customers", customerV2Handler.GetAllCustomers)
	v2.GET("/customer/:id", customerV2Handler.GetCustomerById)
	v2.PATCH("/customer/:id", customerV2Handler.UpdateCustomerById)
	v2.DELETE("/customers", customerV2Handler.DeleteAllCustomers)
	v2.DELETE("/customer/:id", customerV2Handler.DeleteCustomerById)

	// Booking route
	v2.POST("/book", bookingV2Handler.CreateBooking)
	v2.GET("/books", bookingV2Handler.GetAllBookings)
	v2.GET("/book/:id", bookingV2Handler.GetBookingById)
	v2.PATCH("/book/:id", bookingV2Handler.UpdateBookingById)
	v2.PATCH("/book/finished/:id", bookingV2Handler.FinishedStatusBooking)
	v2.DELETE("/books", bookingV2Handler.DeleteAllbookings)
	v2.DELETE("/book/:id", bookingV2Handler.DeleteBookingById)

	// Driver Incentive route (ONLY READ because driver-incentives calucalte automatic while booking)
	v2.GET("/drivers-incentives", driverIncentiveV2Handler.GetAllDriverIncentives)
	v2.GET("/driver-incentive/:id", driverIncentiveV2Handler.GetDriverIncentiveById)

	r.Run(":" + port)

}
