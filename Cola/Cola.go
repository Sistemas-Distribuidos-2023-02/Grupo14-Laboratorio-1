package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

const (
	queueName   = "mi_cola" // Nombre de la cola
	rabbitMQURL = "amqp://guest:guest@localhost:5672/"
)

func main() {
	// Establecer conexión con RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("No se pudo conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Crear canal
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("No se pudo abrir el canal: %v", err)
	}
	defer ch.Close()

	// Declara la cola
	_, err = ch.QueueDeclare(
		queueName, // Nombre de la cola
		false,     // Durable
		false,     // Eliminar cuando no se usa
		false,     // Exclusiva
		false,     // No espera por el mensaje si no hay
		nil,       // Argumentos adicionales
	)
	if err != nil {
		log.Fatalf("No se pudo declarar la cola: %v", err)
	}

	// Consumir mensajes de la cola
	msgs, err := ch.Consume(
		queueName, // Nombre de la cola
		"",        // Consumidor
		true,      // Auto-acknowledge
		false,     // Exclusive
		false,     // No local
		false,     // No wait
		nil,       // Argumentos adicionales
	)
	if err != nil {
		log.Fatalf("No se pudo registrar el consumidor en la cola: %v", err)
	}

	// Consumir mensajes de la cola en el servidor central
	go func() {
		for msg := range msgs {
			fmt.Printf("Mensaje recibido en el servidor central: %s\n", msg.Body)
		}
	}()

	// Simulación: cuatro servidores regionales que publican mensajes
	for i := 1; i <= 4; i++ {
		go func(serverID int) {
			for {
				message := fmt.Sprintf("Mensaje desde el servidor regional %d: %s", serverID, time.Now())
				err := ch.Publish(
					"",        // Exchange
					queueName, // Cola
					false,     // Mandatory
					false,     // Immediate
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        []byte(message),
					},
				)
				if err != nil {
					log.Printf("Error al publicar mensaje: %v", err)
				}
				time.Sleep(5 * time.Second)
			}
		}(i)
	}

	// Mantener el programa en funcionamiento
	forever := make(chan bool)
	<-forever
}
