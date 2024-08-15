package rabbitmq

import (
	"booking-service/genproto/bookings"
	"booking-service/internal/service"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

type MsgBroker struct {
	bookingService   *service.BookingService
	createBooking    <-chan amqp.Delivery
	updateBooking    <-chan amqp.Delivery
	cancelBooking    <-chan amqp.Delivery
	channel          *amqp.Channel
	numberOfServices int
	logger           *log.Logger
	wg               *sync.WaitGroup
}

func New(service *service.BookingService,
	channel *amqp.Channel,
	logger *log.Logger,
	createBooking <-chan amqp.Delivery,
	updateBooking <-chan amqp.Delivery,
	cancelBooking <-chan amqp.Delivery,
	wg *sync.WaitGroup,
	numberOfServices int) *MsgBroker {
	return &MsgBroker{
		bookingService:   service,
		channel:          channel,
		createBooking:    createBooking,
		updateBooking:    updateBooking,
		cancelBooking:    cancelBooking,
		logger:           logger,
		wg:               wg,
		numberOfServices: numberOfServices,
	}
}

func (m *MsgBroker) StartToConsume(ctx context.Context, contentType string) {
	m.wg.Add(m.numberOfServices)
	consumerCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go m.consumeMessages(consumerCtx, m.createBooking, m.bookingService.CreateBooking, "create")
	go m.consumeMessages(consumerCtx, m.updateBooking, m.bookingService.UpdateBooking, "update")
	go m.consumeMessages(consumerCtx, m.cancelBooking, m.bookingService.CancelBooking, "cancel")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	m.logger.Println("Shutting down, waiting for consumers to finish")
	cancel()
	m.wg.Wait()
	m.logger.Println("All consumers have stopped")
}

func (m *MsgBroker) consumeMessages(ctx context.Context, messages <-chan amqp.Delivery, serviceFunc interface{}, logPrefix string) {
	defer m.wg.Done()
	for {
		select {
		case val := <-messages:
			var response proto.Message
			var err error

			switch logPrefix {
			case "create":
				var req bookings.NewBooking
				if err := json.Unmarshal(val.Body, &req); err != nil {
					log.Println("ERROR WHILE UNMARSHALING DATA: %s\n", err.Error())
					val.Nack(false, false)
					continue
				}
				response, err = serviceFunc.(func(context.Context, *bookings.NewBooking) (*bookings.CreateResp, error))(ctx, &req)
			case "update":
				var req bookings.NewData
				if err := json.Unmarshal(val.Body, &req); err != nil {
					m.logger.Println("ERROR WHILE UNMARSHALING DATA: %s\n", err.Error())
					val.Nack(false, false)
					continue
				}
				response, err = serviceFunc.(func(context.Context, *bookings.NewData) (*bookings.UpdateResp, error))(ctx, &req)
			case "cancel":
				var req bookings.ID
				if err := json.Unmarshal(val.Body, &req); err != nil {
					m.logger.Println("ERROR WHILE UNMARSHALING DATA: %s\n", err.Error())
					val.Nack(false, false)
					continue
				}
				response, err = serviceFunc.(func(context.Context, *bookings.ID) (*bookings.Void, error))(ctx, &req)
			}

			if err != nil {
				m.logger.Println("Failed in %s: %s\n", logPrefix, err.Error())
				val.Nack(false, false)
				continue
			}

			val.Ack(false)

			_, err = proto.Marshal(response)
			if err != nil {
				m.logger.Println("Failed to marshal response: %s\n", err.Error())
				continue
			}
		case <-ctx.Done():
			m.logger.Println("Context done, stopping %s consumer", logPrefix)
			return
		}
	}
}
