package main

import (
	"gateway/genproto/auth"
	"gateway/genproto/bookings"
	"gateway/genproto/payments"
	"gateway/genproto/providers"
	"gateway/genproto/services"
	"gateway/internal/api"
	"gateway/internal/api/handler"
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	router := gin.Default()

	authConn, err := grpc.NewClient(":50050", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to gRPC server:", err)
		return
	}
	defer authConn.Close()

	bookingConn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to gRPC server:", err)
		return
	}
	defer bookingConn.Close()

	userClient := auth.NewUserManagementServiceClient(authConn)
	authClient := auth.NewAuthServiceClient(authConn)
	provierClient := providers.NewProvidersClient(bookingConn)
	serClient := services.NewServicesClient(bookingConn)
	bookingsClient := bookings.NewBookingsClient(bookingConn)
	paymentClient := payments.NewPaymentsClient(bookingConn)

	userHandler := handler.NewUserManagementHandler(userClient)
	authHandler := handler.NewAuthHandler(authClient, userClient)
	proHandler := handler.NewProviderManagementHandler(provierClient)
	serHandler := handler.NewServiceManagementHandler(serClient)
	bookingHandler := handler.NewBookingHandler(bookingsClient)
	paymentsHandler := handler.NewPaymentHandler(paymentClient)

	mainHandler := handler.NewMainHandler(
		userHandler,
		authHandler,
		proHandler,
		serHandler,
		bookingHandler,
		paymentsHandler,
	)

	api.SetupRouter(router, mainHandler)

	router.Run(":8080")
}
