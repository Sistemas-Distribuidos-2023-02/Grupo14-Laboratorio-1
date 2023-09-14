package main

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	pb "github.com/Sistemas-Distribuidos-2023-02/Grupo14-Laboratorio-1/proto" // Asegúrate de usar la ruta correcta a tus archivos .proto
)

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

type regionalServer struct{}

func (s *regionalServer) ReceiveMessage(ctx context.Context, req *pb.Message) (*pb.Response, error) {
	// Manejar el mensaje recibido desde el servidor central
	//content := req.Content

	// Realizar cualquier lógica necesaria

	// Responder al servidor central
	response := &pb.Response{Message: "Mensaje recibido en el servidor regional"}
	return response, nil
}

func main() {
	keys := generateRandomNum()
	fmt.Printf("Se pediran %d claves", keys)
	/* listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("No se pudo escuchar en el puerto 50051: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterRegionalServerServer(server, &regionalServer{})

	fmt.Println("Servidor Regional gRPC escuchando en el puerto 50051...")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	} */
}
