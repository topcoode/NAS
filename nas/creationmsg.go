package nas

import (
	"encoding/binary"
	"fmt"
)

type NASMessage struct {
	MessageType uint8
	Sequence    uint16
	Payload     []byte
}

func CreationMessage() {
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

	// Print the packed message data
	fmt.Printf("Packed NAS message data: %X\n", messageData)
}
