package main

import (
	"fmt"
	"nas/nas"
	"net"
)

func main() {
	nas.CreationMessage()
	// Create a TCP listener
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Failed to create listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on", listener.Addr())

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle each connection concurrently
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Your code to handle the NAS protocol goes here

	// Example: Read data from the client
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	// Process the received data
	data := buffer[:n]
	fmt.Println("Received data:", string(data))

	// Example: Send a response to the client
	response := []byte("Hello from NAS server!")
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error sending response:", err)
		return
	}
}
