package main

import (
	"booking-service/config"
	"booking-service/genproto/bookings"
	"booking-service/genproto/payments"
	"booking-service/genproto/providers"
	"booking-service/genproto/reviews"
	"booking-service/genproto/services"
	"booking-service/internal/rabbitmq"
	"booking-service/internal/service"
	mongostorage "booking-service/internal/storage/mongo"
	mongodb "booking-service/pkg/mongo"
	"context"
	"log"
	"net"
	"os"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

func main() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	cnf := config.NewConfig()
	if err := cnf.Load(); err != nil {
		logger.Fatalf("Error loading config: %v", err)
		return
	}

	lis, err := net.Listen("tcp", cnf.Server.Host+":"+cnf.Server.Port)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()

	mongoClient, err := mongodb.ConnectDB(cnf)
	if err != nil {
		logger.Fatal(err)
	}

	providerMongoStorage := mongostorage.NewProviderMongoStorage(mongoClient.Collection("providers"))
	serviceStorage := mongostorage.NewServicesMongoStorage(mongoClient.Collection("services"))
	bookingsStorage := mongostorage.NewBookingsStorage(mongoClient.Collection("bookings"))
	paymentStorage := mongostorage.NewPaymentsStorage(mongoClient.Collection("payments"))
	reviewsStorage := mongostorage.NewReviewsStorage(mongoClient.Collection("reviews"))

	providerService := service.NewProviderService(logger, providerMongoStorage)
	servicesService := service.NewServicesService(logger, serviceStorage)
	bookingsService := service.NewBookingService(logger, bookingsStorage)
	paymentsService := service.NewPaymentsService(logger, paymentStorage)
	reviewsService := service.NewReviewService(logger, reviewsStorage)

	bookings.RegisterBookingsServer(grpcServer, bookingsService)
	services.RegisterServicesServer(grpcServer, servicesService)
	providers.RegisterProvidersServer(grpcServer, providerService)
	payments.RegisterPaymentsServer(grpcServer, paymentsService)
	reviews.RegisterReviewsServer(grpcServer, reviewsService)

	logger.Println("Serving on port 50051")

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	var conn *amqp.Connection

	time.Sleep(time.Second * 15)
	conn, err = amqp.Dial(cnf.RabbitMQ.RabbitMQ)
	if err != nil {
		logger.Fatal("error connecting to RabbitMQ: ", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	createQ, err := getQueue(ch, "create")
	if err != nil {
		logger.Fatalf("Failed to declare 'create' queue: %v", err)
	}

	updateQ, err := getQueue(ch, "update")
	if err != nil {
		logger.Fatalf("Failed to declare 'update' queue: %v", err)
	}

	cancelQ, err := getQueue(ch, "cancel")
	if err != nil {
		logger.Fatalf("Failed to declare 'cancel' queue: %v", err)
	}

	createBooking, err := getMessageQueue(ch, createQ)
	if err != nil {
		logger.Fatalf("Failed to consume from 'create' queue: %v", err)
	}

	updateBooking, err := getMessageQueue(ch, updateQ)
	if err != nil {
		logger.Fatalf("Failed to consume from 'update' queue: %v", err)
	}

	cancelBooking, err := getMessageQueue(ch, cancelQ)
	if err != nil {
		logger.Fatalf("Failed to consume from 'cancel' queue: %v", err)
	}

	wg := sync.WaitGroup{}
	msgBroker := rabbitmq.New(bookingsService, ch, logger, createBooking, updateBooking, cancelBooking, &wg, 3)

	msgBroker.StartToConsume(context.Background(), "application/json")
}

func getQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
}

func getMessageQueue(ch *amqp.Channel, q amqp.Queue) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}
