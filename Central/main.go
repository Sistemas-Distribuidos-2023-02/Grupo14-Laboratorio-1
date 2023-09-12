package main

import (
	"context"
	"fmt"
	"log"

	pb "Lab1SD/proto" // Asegúrate de usar la ruta correcta a tus archivos .proto

	"google.golang.org/grpc"
)

func main() {
	// Establecer conexión gRPC con el servidor regional
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("No se pudo conectar al servidor regional: %v", err)
	}
	defer conn.Close()
	client := pb.NewRegionalServerClient(conn)

	// Enviar mensaje al servidor regional
	message := &pb.Message{Content: "Mensaje desde el servidor central"}
	response, err := client.ReceiveMessage(context.Background(), message)
	if err != nil {
		log.Fatalf("Fallo al enviar mensaje al servidor regional: %v", err)
	}

	// Procesar la respuesta del servidor regional
	fmt.Printf("Respuesta del servidor regional: %s\n", response.Message)
}
