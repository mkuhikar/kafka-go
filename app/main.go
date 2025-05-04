package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	buffer := make([]byte, 512)

	n, err := conn.Read(buffer) //The buffer is a pre-allocated memory slice (512 bytes in this case) used to store incoming data from the client.
	if err != nil {
		fmt.Println("Error reading request:", err.Error())
		os.Exit(1)
	}
	fmt.Printf("the no of bytes received from client are %d", n) /*fmt.Printf allows formatted strings, just like printf in C.

	%d is a placeholder for an integer — here, it's replaced with the value of n.

	\n adds a newline after the output. */

	// correlation_id := binary.BigEndian.Uint32(buffer[4:8]) //converted 4 bytes of buffer to uInt32, in big endian most significant byte comes first
	// message_size := uint32(26 - 4)
	// api_version := binary.BigEndian.Uint16(buffer[2:4]) //It should be the number of bytes in your response body + header, excluding the first 4 bytes (the message_size field itself).
	// Prepare values
	correlation_id := binary.BigEndian.Uint32(buffer[4:8])
	api_version := binary.BigEndian.Uint16(buffer[2:4])
	code := uint16(0)
	if api_version > 4 {
		code = 35
	}

	// Use a buffer to construct the response dynamically
	var response []byte

	// Reserve space for message size (4 bytes, to be filled later)
	response = append(response, 0, 0, 0, 0)

	// Correlation ID (4 bytes)
	tmp := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, correlation_id)
	response = append(response, tmp...)

	// Error Code (2 bytes)
	tmp = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, code)
	response = append(response, tmp...)

	// Number of APIs (1 byte)
	response = append(response, 2)

	// First API Key (18), MinVersion (3), MaxVersion (4)
	tmp = make([]byte, 6)
	binary.BigEndian.PutUint16(tmp[0:2], 18)
	binary.BigEndian.PutUint16(tmp[2:4], 3)
	binary.BigEndian.PutUint16(tmp[4:6], 4)
	response = append(response, tmp...)

	// Tagged fields for API list (1 byte)
	response = append(response, 0)

	// Throttle time (4 bytes)
	tmp = make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, 0)
	response = append(response, tmp...)

	// Tagged fields for throttle time (1 byte)
	response = append(response, 0)

	// Set message size at the beginning (excluding the 4 bytes for size itself)
	binary.BigEndian.PutUint32(response[0:4], uint32(len(response)-4))
	// response:= []byte{0,0,0,0,0,0,0,7} // just a hard coded way to send correlation id
	fmt.Println(response)

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)

	}
	// fmt.Println("Here ", l)

}
