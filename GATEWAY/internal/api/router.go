package api

import (
	"gateway/internal/api/handler"

	_ "gateway/cmd/docs"

	"github.com/gin-gonic/gin"
	files "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Package api API.
//
// @title # Car Wash System
// @version 1.0
// @description API Endpoints for car wash app
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func SetupRouter(r *gin.Engine, handler handler.MainHandler) {
	r.GET("swagger/*any", ginSwagger.WrapHandler(files.Handler))
	auth := r.Group("/auth")
	{
		auth.POST("/register", handler.Authentication().Register)
		auth.POST("/login", handler.Authentication().Login)
	}

	user := r.Group("/user")
	{
		user.POST("/", handler.UserManagement().CreateUser)
		user.GET("/:id", handler.UserManagement().GetUserByID)
		user.PUT("/", handler.UserManagement().UpdateUser)
		user.DELETE("/:id", handler.UserManagement().DeleteUser)
		user.GET("/", handler.UserManagement().GetUsers)
	}

	providers := r.Group("/provider")
	{
		providers.POST("/register", handler.ProviderManagement().RegisterProvider)
		providers.POST("/search", handler.ProviderManagement().SearchProviders)
	}
	services := r.Group("/services")
	{
		services.POST("/", handler.ServiceManagement().CreateService)
		services.GET("/", handler.ServiceManagement().ListServices)
		services.GET("/:id", handler.ServiceManagement().GetServiceByID)
		services.PUT("/:id", handler.ServiceManagement().UpdateService)
		services.DELETE("/:id", handler.ServiceManagement().DeleteService)
	}

	bookings := r.Group("/booking")
	{
		bookings.POST("", handler.BookingsManagement().CreateBooking)
		bookings.POST("/search", handler.BookingsManagement().ListBookings)
		bookings.GET("/:id", handler.BookingsManagement().GetBooking)
		bookings.PUT("/:id", handler.BookingsManagement().UpdateBooking)
		bookings.DELETE("/:id", handler.BookingsManagement().DeleteBooking)
	}

	payments := r.Group("/payments")
	{
		payments.POST("", handler.PaymentManagement().CreatePayment)
		payments.GET("/:id", handler.PaymentManagement().GetPayment)
		payments.GET("", handler.PaymentManagement().ListPayments)
	}

	reviews := r.Group("/reviews")
	{
		reviews.POST("/", handler.ReviewManagement().CreateReview)
		reviews.GET("/", handler.ReviewManagement().ListReviews)
		// reviews.GET("/:id", handler.ReviewManagement().GetReview)
		reviews.PUT("/:id", handler.ReviewManagement().UpdateReview)
		reviews.DELETE("/:id", handler.ReviewManagement().DeleteReview)
	}
}
