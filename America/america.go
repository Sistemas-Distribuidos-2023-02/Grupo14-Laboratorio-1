package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "Lab1SD/proto" // Asegúrate de usar la ruta correcta a tus archivos .proto

	"google.golang.org/grpc"
)

type regionalServer struct{}

func (s *regionalServer) ReceiveMessage(ctx context.Context, req *pb.Message) (*pb.Response, error) {
	// Manejar el mensaje recibido desde el servidor central
	content := req.Content

	// Realizar cualquier lógica necesaria

	// Responder al servidor central
	response := &pb.Response{Message: "Mensaje recibido en el servidor regional"}
	return response, nil
}

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("No se pudo escuchar en el puerto 50051: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterRegionalServerServer(server, &regionalServer{})

	fmt.Println("Servidor Regional gRPC escuchando en el puerto 50051...")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
