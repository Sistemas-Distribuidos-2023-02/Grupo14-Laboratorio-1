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
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

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
	// Generar un número aleatorio entre min y max (inclusive).
	return rand.Intn(max-min+1) + min
}

// Función para enviar un valor numérico a los servidores regionales
func sendValueToRegionals(val int) {
	// Direcciones de los servidores regionales
	serverAddresses := []string{"localhost:50051", "localhost:50052", "localhost:50053"} // Agrega las direcciones de tus servidores regionales

	// Valor numérico a enviar
	value := val
	msg := strconv.Itoa(value)

	for _, address := range serverAddresses {
		// Establecer conexión gRPC con el servidor regional
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Printf("No se pudo conectar al servidor regional en %s: %v", address, err)
			continue
		}
		defer conn.Close()

		client := pb.NewRegionalServerClient(conn) // puede que este alverrre

		// Crear un contexto con un límite de tiempo para la llamada gRPC
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Enviar el valor numérico al servidor regional
		response, err := client.ReceiveMessage(ctx, &pb.Message{Content: msg})
		if err != nil {
			log.Printf("Fallo al enviar el valor al servidor regional en %s: %v", address, err)
			continue
		}

		// Procesar la respuesta del servidor regional
		log.Printf("Respuesta del servidor regional en %s: %s", address, response.Message)
	}
}

func main() {

	numeros := archivo()
	//fmt.Println("los numeros son:", numeros[0], numeros[1], numeros[2])

	iterations := numeros[2]
	for i := 0; i < iterations || iterations == -1; i++ {
		fmt.Printf("Generacion %d/%d\n", i+1, iterations)
		numKeys := generateRandomKeys(numeros[0], numeros[1])
		fmt.Printf("Generando %d llaves de acceso...\n", numKeys)
		sendValueToRegionals(numKeys)
	}

}
