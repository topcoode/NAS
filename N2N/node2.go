package main

import (
	"fmt"
	"net"
)

func main() {
	listenAddr := "127.0.0.1:8081" // Example listening address and port

	// Listen for incoming connections
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Waiting for NAS message...")

	// Accept incoming connections
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	// Read the incoming NAS message
	buffer := make([]byte, 1024) // Adjust the buffer size based on your message size
	bytesRead, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading NAS message:", err)
		return
	}

	receivedData := buffer[:bytesRead]

	// Process the received NAS message
	fmt.Printf("Received NAS message: %X\n", receivedData)
}
