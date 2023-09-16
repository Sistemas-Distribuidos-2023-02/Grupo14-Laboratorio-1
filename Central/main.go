package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo14-Laboratorio-1/proto" // Asegúrate de usar la ruta correcta a tus archivos .proto
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup // no se usarlo??
var numKeys int

func archivo() [3]int {
	file, err := os.Open("Central/parametros_de_inicio.txt")
	if err != nil {
		fmt.Println("File reading error", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	linea1 := scanner.Text()
	partes := strings.Split(linea1, "-")
	numero1, err1 := strconv.Atoi(partes[0])
	numero2, err2 := strconv.Atoi(partes[1])

	if err1 != nil || err2 != nil {
		fmt.Println("Error al convertir los números:", err1, err2)
	}

	min := numero1
	max := numero2

	scanner.Scan()
	linea2 := scanner.Text()
	numero3, err3 := strconv.Atoi(linea2)

	if err3 != nil {
		fmt.Println("Error al convertir el número:", err3)
	}

	iter := numero3

	var lista [3]int
	lista[0] = min
	lista[1] = max
	lista[2] = iter
	return lista
}

func generateRandomKeys(min, max int) int {
	return rand.Intn(max-min+1) + min
}

// Función para enviar un valor numérico a los servidores regionales
func sendValueToRegionals(val string) {
	// Direcciones de los servidores regionales
	serverAddresses := []string{"localhost:50051-america"}
	//serverAddresses := {"dist053.inf.santiago.usm.cl:50051-america", "dist054.inf.santiago.usm.cl:50051-asia", "dist055.inf.santiago.usm.cl:50051-europa", "dist056.inf.santiago.usm.cl:50051-oceania"}

	// Valor numérico a enviar
	msg := val
	for _, address := range serverAddresses {
		// Establecer conexión gRPC con el servidor regional
		partes := strings.Split(string(address), "-")
		address = partes[0]
		nombre := partes[1]
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Printf("No se pudo conectar al servidor regional en %s: %v", address, err)
			continue
		}
		defer conn.Close()

		client := pb.NewRegionalServerClient(conn) // puede que este alverrre

		// Crear un contexto con un límite de tiempo para la llamada gRPC
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Enviar el valor numérico al servidor regional
		response, err := client.ReceiveMessage(ctx, &pb.Message{Content: msg, Name: nombre})
		if err != nil {
			log.Printf("Fallo al enviar el valor al servidor regional en %s: %v", address, err)
			continue
		}

		// Procesar la respuesta del servidor regional
		log.Printf("Respuesta del servidor %s en %s: %s", nombre, address, response.Message)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func consumeMessages() {
	// Establece la conexión a RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Crea un canal
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declara la cola desde la que se consumirán los mensajes
	queueName := "mi_cola"
	_, err = ch.QueueDeclare(
		queueName, // Nombre de la cola
		true,      // Duradera
		false,     // Eliminable cuando no se use
		false,     // Exclusiva (solo se puede usar desde esta conexión)
		false,     // No se requieren argumentos adicionales
		nil,       // Argumentos adicionales
	)
	failOnError(err, "Failed to declare a queue")

	// Configura el consumidor
	msgs, err := ch.Consume(
		queueName, // Nombre de la cola
		"",        // Etiqueta del consumidor (deja en blanco para una etiqueta generada automáticamente)
		true,      // Autoacknowledgment (RabbitMQ marcará automáticamente los mensajes como entregados)
		false,     // Exclusividad del consumidor (puede haber múltiples consumidores en la misma cola)
		false,     // NoLocal (no se entregan mensajes publicados por el propio consumidor)
		false,     // NoWait (no esperar a que se confirme la solicitud)
		nil,       // Argumentos adicionales
	)
	failOnError(err, "Failed to register a consumer")

	// Escucha los mensajes
	for msg := range msgs {
		time.Sleep(5 * time.Second)
		log.Printf("Received a message: %s", msg.Body)
		partes := strings.Split(string(msg.Body), " ")
		numero, err1 := strconv.Atoi(partes[0])

		if err1 != nil {
			fmt.Println("Error al convertir los números:", err1)
		}

		if numKeys-numero > 0 {
			numKeys = numKeys - numero
			sendValueToRegionals((strconv.Itoa(numero) + " " + "asignadas"))
		} else {
			sendValueToRegionals((strconv.Itoa(numKeys) + " " + "asignadas"))
			numKeys = 0
		}

	}
}

func main() {

	numeros := archivo()
	//fmt.Println("los numeros son:", numeros[0], numeros[1], numeros[2])

	iterations := numeros[2]
	for i := 0; i < iterations || iterations == -1; i++ {
		if iterations != -1 {
			fmt.Printf("Generacion %d/%d\n", i+1, iterations)
		} else {
			fmt.Printf("Generacion %d/Infinito\n", i+1)
		}

		numKeys = generateRandomKeys(numeros[0], numeros[1])
		fmt.Printf("Generando %d llaves de acceso...\n", numKeys)
		sendValueToRegionals((strconv.Itoa(numKeys) + " " + "disponibles"))

		consumeMessages()

	}

}
