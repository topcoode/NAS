package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

type NASMessage struct {
	MessageType uint8
	Sequence    uint16
	Payload     []byte
}

func main() {
	// Create a new NAS message
	nasMsg := NASMessage{
		MessageType: 0x01,             // Example message type
		Sequence:    0x1234,           // Example sequence number
		Payload:     []byte("Hello!"), // Example payload
	}

	// Pack the message fields into a byte array
	messageData := make([]byte, 0)
	messageData = append(messageData, nasMsg.MessageType)
	sequenceBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(sequenceBytes, nasMsg.Sequence)
	messageData = append(messageData, sequenceBytes...)
	messageData = append(messageData, nasMsg.Payload...)

	// Send the message to the destination node
	destAddr := "127.0.0.1:8081" // Example destination address and port
	conn, err := net.Dial("tcp", destAddr)
	if err != nil {
		fmt.Println("Error connecting to the destination node:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(messageData)
	if err != nil {
		fmt.Println("Error sending NAS message:", err)
		return
	}

	fmt.Println("NAS message sent successfully.")
}
