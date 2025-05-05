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

	correlation_id := binary.BigEndian.Uint32(buffer[4:8]) //converted 4 bytes of buffer to uInt32, in big endian most significant byte comes first
	// message_size := uint32(26 - 4)                          //It should be the number of bytes in your response body + header, excluding the first 4 bytes (the message_size field itself).
	// response := make([]byte, 16)
	// binary.BigEndian.PutUint32(response[0:4], message_size)
	// binary.BigEndian.PutUint32(response[4:8], correlation_id) //pass uint32 correlation id , correlation id consits of 4 byte, each byte is of 2 hex digits
	// binary.BigEndian.PutUint16(response[8:10], 0)
	// response = append(response, 1) // num_apis = 1
	// tmp := make([]byte, 6)
	// binary.BigEndian.PutUint16(tmp[0:2], 18) // api key
	// binary.BigEndian.PutUint16(tmp[2:4], 0)  // min version
	// binary.BigEndian.PutUint16(tmp[4:6], 4)  // max version
	// response = append(response, tmp...)
	// response = append(response, 0x00)                       // tagged fields
	// response = append(response, 0x00, 0x00, 0x00, 0x00)     // throttle time
	// response = append(response, 0x00)                       // tagged fields
	// message_size = uint32(len(response))                    // The actual length of the response buffer excluding the first 4 bytes
	// binary.BigEndian.PutUint32(response[0:4], message_size) // Update message_size with the correct total length
	// response:= []byte{0,0,0,0,0,0,0,7} // just a hard coded way to send correlation id
	var response []byte

	response = append(response, 0, 0, 0, 0) // placeholder for message size

	tmp := make([]byte, 4)
	binary.BigEndian.PutUint32(tmp, correlation_id)
	response = append(response, tmp...) // correlation_id

	tmp = make([]byte, 2)
	binary.BigEndian.PutUint16(tmp, 0)
	response = append(response, tmp...) // error code = 0

	response = append(response, 1) // num_apis = 1

	tmp = make([]byte, 6)
	binary.BigEndian.PutUint16(tmp[0:2], 18) // api key
	binary.BigEndian.PutUint16(tmp[2:4], 0)  // min version
	binary.BigEndian.PutUint16(tmp[4:6], 4)  // max version
	response = append(response, tmp...)

	response = append(response, 0x00)                   // tagged fields
	response = append(response, 0x00, 0x00, 0x00, 0x00) // throttle time
	response = append(response, 0x00)                   // tagged fields

	// Now fix the message size (excluding the 4 bytes for the size itself)
	binary.BigEndian.PutUint32(response[0:4], uint32(len(response)-4))
	fmt.Println(response)

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)

	}
	// fmt.Println("Here ", l)

}
