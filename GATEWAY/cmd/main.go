package main

import (
	"context"
	"gateway/config"
	"gateway/genproto/auth"
	"gateway/genproto/bookings"
	"gateway/genproto/payments"
	"gateway/genproto/providers"
	"gateway/genproto/reviews"
	"gateway/genproto/services"
	"gateway/internal/api"
	"gateway/internal/api/handler"
	"gateway/internal/msgbroker"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cnf := config.NewConfig()
	if err := cnf.Load(); err != nil {
		log.Fatal(err)
		return
	}

	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	router := gin.Default()

	authConn, err := grpc.NewClient("auth:"+cnf.Auth, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to gRPC server:", err)
		return
	}
	defer authConn.Close()

	bookingConn, err := grpc.NewClient("booking:"+cnf.Booking, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect to gRPC server:", err)
		return
	}
	log.Println()
	defer bookingConn.Close()

	userClient := auth.NewUserManagementServiceClient(authConn)
	authClient := auth.NewAuthServiceClient(authConn)
	provierClient := providers.NewProvidersClient(bookingConn)
	serClient := services.NewServicesClient(bookingConn)
	bookingsClient := bookings.NewBookingsClient(bookingConn)
	paymentClient := payments.NewPaymentsClient(bookingConn)
	reviewsClient := reviews.NewReviewsClient(bookingConn)

	userHandler := handler.NewUserManagementHandler(userClient)
	authHandler := handler.NewAuthHandler(authClient, userClient)
	proHandler := handler.NewProviderManagementHandler(provierClient)
	serHandler := handler.NewServiceManagementHandler(serClient)
	paymentsHandler := handler.NewPaymentHandler(paymentClient)
	reviewsHandler := handler.NewReviewHandler(reviewsClient)
	var conn *amqp.Connection

	time.Sleep(time.Second * 15)
	conn, err = amqp.Dial(cnf.RabbitMQ.RabbitMQ)
	if err != nil {
		logger.Fatal("error connecting to RabbitMQ: ", err)
	}

	if conn != nil {
		logger.Println("Connection created")
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("Error opening channel: %v", err)
	}
	defer ch.Close()

	createBooking, createq, err := getMessages("create", ch)
	if err != nil {
		logger.Fatalf("Error getting registration messages: %v", err)
	}

	updateBooking, updateq, err := getMessages("update", ch)
	if err != nil {
		logger.Fatalf("Error getting update messages: %v", err)
	}

	cancelBooking, cancelq, err := getMessages("delete", ch)
	if err != nil {
		logger.Fatalf("Error getting deletion messages: %v", err)
	}

	msgBroker, err := msgbroker.NewMsgBorker(ch, 10*time.Second, createBooking, updateBooking, cancelBooking, context.Background())
	if err != nil {
		logger.Fatalf("Error creating RPC client: %v", err)
	}

	bookingHandler := handler.NewBookingHandler(bookingsClient, createq, updateq, cancelq, msgBroker, logger)

	mainHandler := handler.NewMainHandler(
		userHandler,
		authHandler,
		proHandler,
		serHandler,
		bookingHandler,
		paymentsHandler,
		reviewsHandler,
	)

	api.SetupRouter(router, mainHandler)

	router.Run(cnf.Server.Host + ":" + cnf.Server.Port)

}

func getMessages(queueName string, ch *amqp.Channel) (<-chan amqp.Delivery, amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, q, err
	}

	messages, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	return messages, q, err
}
