package msgbroker

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MsgBroker struct {
	Ch            *amqp.Channel
	timeoutCh     <-chan time.Time
	CreateBooking <-chan amqp.Delivery
	UpdateBooking <-chan amqp.Delivery
	CancelBooking <-chan amqp.Delivery
	Ctx           context.Context
}

func NewMsgBorker(ch *amqp.Channel, timeout time.Duration, createBooking <-chan amqp.Delivery, updateBooking <-chan amqp.Delivery, cancelBooking <-chan amqp.Delivery, ctx context.Context) (*MsgBroker, error) {
	timeoutCh := time.After(timeout)
	return &MsgBroker{
		Ch:            ch,
		timeoutCh:     timeoutCh,
		CreateBooking: createBooking,
		UpdateBooking: updateBooking,
		CancelBooking: cancelBooking,
		Ctx:           ctx,
	}, nil
}

func (m *MsgBroker) PublishToQueue(messages <-chan amqp.Delivery, body []byte, q amqp.Queue, contentType string) error {
	corrId := uuid.New().String()
	fmt.Println("Generated a new CorrelationId:", corrId)

	return m.Ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:   contentType,
			CorrelationId: corrId,
			Body:          body,
		},
	)
}
