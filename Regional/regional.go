package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"google.golang.org/grpc"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo14-Laboratorio-1/proto"
	"github.com/streadway/amqp"
)

var request int
var serverActive bool = true

func generateRandomNum() int {
	file, err := os.Open("Regional/parametros_de_inicio.txt")
	if err != nil {
		fmt.Println("File reading error", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	linea := scanner.Text()
	numero1, err1 := strconv.Atoi(linea)
	if err1 != nil {
		fmt.Println("Error al convertir el número:", err1)
	}

	keys := numero1
	max := keys/2 + keys*10/100
	min := keys/2 - keys*10/100
	pedir := rand.Intn(max-min+1) + min
	return pedir
}

func ConnectToRabbitMQ(url, queueName string) error {
	// Conectarse a RabbitMQ
	connection, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer connection.Close()

	// Crear un canal de comunicación
	channel, err := connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// Declarar la cola a la que te quieres conectar
	_, err = channel.QueueDeclare(
		queueName, // Nombre de la cola
		true,      // Durable
		false,     // Autoeliminable
		false,     // Exclusiva
		false,     // No esperar a que se confirme la entrega
		nil,       // Argumentos adicionales
	)
	if err != nil {
		return err
	}

	// Puedes realizar operaciones adicionales con la cola aquí si es necesario.

	return nil
}

type RegionalServer struct {
	pb.UnimplementedRegionalServerServer
}

// ReceiveMessage es la implementación del método ReceiveMessage
func (s *RegionalServer) ReceiveMessage(ctx context.Context, req *pb.Message) (*pb.Response, error) {
	content := req.Content
	fmt.Printf("Mensaje recibido: %s\n", content)
	partes := strings.Split(string(content), " ")
	keys, err1 := strconv.Atoi(partes[0])
	tipo := partes[1]
	if err1 != nil {
		log.Fatalf("Error al escuchar: %v", err1)
	}
	if tipo == "disponibles" {
		fmt.Printf("Pido: %d\n", request)
	}

	if tipo == "asignadas" {
		if request <= keys {
			request = 0
			fmt.Printf("me faltan: %d\n", request)
			serverActive = false
		} else {
			request = request - keys // no es correcto pero casi
			fmt.Printf("me faltan: %d\n", keys)
		}
	}

	// Envia el mensaje a RabbitMQ
	err := PublishToRabbitMQ((strconv.Itoa(request) + " " + "nombre-server"))
	if err != nil {
		log.Fatalf("Error al enviar el mensaje a RabbitMQ: %v", err)
	}
	return &pb.Response{Message: "Mensaje recibido"}, nil
}

func PublishToRabbitMQ(message string) error {
	// Configura la conexión a RabbitMQ
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer connection.Close()

	// Crea un canal
	channel, err := connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	// Declara la cola a la que deseas enviar el mensaje
	queueName := "mi_cola"
	_, err = channel.QueueDeclare(
		queueName, // Nombre de la cola
		true,      // Durable
		false,     // Autoeliminable
		false,     // Exclusiva
		false,     // No esperar a que se confirme la entrega
		nil,       // Argumentos adicionales
	)
	if err != nil {
		return err
	}

	// Publica el mensaje en la cola
	err = channel.Publish(
		"",        // Intercambio
		queueName, // Cola
		false,     // Mandatorio
		false,     // Publicación inmediata
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

func shouldServerStop() bool {
	return !serverActive && request == 0
}

func main() {
	// Crea un servidor gRPC y registra el servicio
	request = generateRandomNum()
	lis, err := net.Listen("tcp", ":50051") // Escucha en el puerto 50051
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRegionalServerServer(grpcServer, &RegionalServer{})

	// Configura una señal para cerrar el servidor cuando se recibe SIGINT (Ctrl+C) o SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Inicia el servidor en una goroutine
	go func() {
		fmt.Println("Servidor gRPC escuchando en :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Error al servir: %v", err)
		}
	}()

	// Espera a que se reciba la señal para cerrar el servidor
	<-stop
	fmt.Println("Cerrando el servidor gRPC...")

	// Verifica si el servidor debe detenerse después de recibir la señal
	if shouldServerStop() {
		fmt.Println("Deteniendo el servidor gRPC...")
		grpcServer.GracefulStop()
		fmt.Println("Servidor gRPC cerrado.")
	}
}
