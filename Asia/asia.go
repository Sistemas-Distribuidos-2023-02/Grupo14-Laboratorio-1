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
	serverID    = 2 // Identificador único del servidor regional
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

	// Simulación: el servidor regional publica mensajes en la cola
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
}
