package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher интерфейс для отправки сообщений в очередь.
type Publisher interface {
	Publish(ctx context.Context, msg interface{}) error
	Close() error
}

// Consumer интерфейс для получения сообщений из очереди.
type Consumer interface {
	Consume() (<-chan amqp.Delivery, error)
	Close() error
}

// Client - реализация клиента RabbitMQ.
type Client struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}

func NewClient(user, password, host string, port int, queueName string) (*Client, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", user, password, host, port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &Client{
		conn:  conn,
		ch:    ch,
		queue: queueName,
	}, nil
}

func (c *Client) Publish(ctx context.Context, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return c.ch.PublishWithContext(ctx,
		"",      // exchange
		c.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (c *Client) Consume() (<-chan amqp.Delivery, error) {
	return c.ch.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack (setting to false for reliability)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
}

func (c *Client) Close() error {
	if c.ch != nil {
		_ = c.ch.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
