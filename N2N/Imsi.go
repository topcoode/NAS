package main

import (
	"encoding/binary"
	"fmt"
)

type AttachRequest struct {
	Header                uint8
	Length                uint16
	SecurityHeader        uint8
	ProtocolDiscriminator uint8
	IMSI                  string
	// Define other fields required by the Attach Request message
}

func main() {
	// Create a new Attach Request message with IMSI
	attachReq := AttachRequest{
		Header:                0x48,              // Example Header value
		Length:                0x0012,            // Example Length value
		SecurityHeader:        0x0E,              // Example Security Header value
		ProtocolDiscriminator: 0x7E,              // Example Protocol Discriminator value
		IMSI:                  "123456789012345", // Example IMSI value
		// Assign values to other fields of the Attach Request message
	}

	// Pack the Attach Request fields into a byte array
	messageData := make([]byte, 0)
	messageData = append(messageData, attachReq.Header)
	lengthBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lengthBytes, attachReq.Length)
	messageData = append(messageData, lengthBytes...)
	messageData = append(messageData, attachReq.SecurityHeader)
	messageData = append(messageData, attachReq.ProtocolDiscriminator)

	imsiBytes := []byte(attachReq.IMSI)
	messageData = append(messageData, imsiBytes...)
	// Pack and append other fields to the messageData byte array

	// Print the packed Attach Request message data
	fmt.Printf("Packed Attach Request message data: %X\n", messageData)
}
