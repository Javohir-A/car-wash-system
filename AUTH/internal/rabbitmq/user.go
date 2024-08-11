package rabbitmq

import (
	"auth-service/genprotos/auth"
	"auth-service/internal/storage"
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type User struct {
	ch      *amqp.Channel
	storage storage.Storage
	queue   string
}

// NewUser initializes the User struct with a RabbitMQ channel and a queue name.
func NewUser(ch *amqp.Channel, queue string, storage storage.Storage) (*User, error) {
	u := &User{
		ch:      ch,
		queue:   queue,
		storage: storage,
	}

	// Ensure the queue exists, and bind it to an exchange if necessary.
	_, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// CreateUser publishes a user creation message to RabbitMQ.
func (u *User) CreateUser(ctx context.Context, user *auth.UserResponse) error {
	// Serialize the user data to JSON.
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Publish the message to the queue.
	err = u.ch.Publish(
		"",      // exchange
		u.queue, // routing key (queue name)
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// Consume starts consuming messages from the RabbitMQ queue.
func (u *User) Consume() {
	msgs, err := u.ch.Consume(
		u.queue, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Process messages in a goroutine.
	go func() {
		for msg := range msgs {
			// Process each message.
			var user auth.UserResponse
			err := json.Unmarshal(msg.Body, &user)
			if err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			// Implement business logic here.
			// err = u.storage.UserManagement().CreateUser(user)
			// if err != nil {
			// 	log.Printf("Failed to save user: %v", err)
			// } else {
			// 	log.Printf("User saved: %s", user.Id)
			// }
		}
	}()
}
