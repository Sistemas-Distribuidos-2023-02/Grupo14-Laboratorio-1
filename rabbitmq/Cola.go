package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

// Config contiene la configuración de RabbitMQ
type Config struct {
	URL   string // URL de conexión a RabbitMQ, por ejemplo, "amqp://username:password@localhost:5672/"
	Queue string // Nombre de la cola
}

// NewConnection crea y devuelve una conexión a RabbitMQ
func NewConnection(cfg Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// PublishMessage publica un mensaje en la cola RabbitMQ
func PublishMessage(conn *amqp.Connection, cfg Config, message string) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(
		cfg.Queue, // Nombre de la cola
		false,     // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",        // exchange
		cfg.Queue, // routingKey
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// ConsumeMessages consume mensajes de la cola RabbitMQ
func ConsumeMessages(conn *amqp.Connection, cfg Config, handler func(message string)) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	q, err := channel.QueueDeclare(
		cfg.Queue, // Nombre de la cola
		false,     // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handler(string(d.Body))
		}
	}()

	log.Printf("Esperando mensajes. Para salir, presione CTRL+C")
	<-forever

	return nil
}
