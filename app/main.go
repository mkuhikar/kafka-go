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
	message_size := uint32(26 - 4)
	api_version := binary.BigEndian.Uint16(buffer[2:4]) //It should be the number of bytes in your response body + header, excluding the first 4 bytes (the message_size field itself).
	response := make([]byte, 16)
	code := uint16(0)
	if api_version > 4 {
		code = uint16(35)
	}
	// binary.BigEndian.PutUint32(response[0:4], message_size)
	binary.BigEndian.PutUint32(response[4:8], correlation_id) //pass uint32 correlation id , correlation id consits of 4 byte, each byte is of 2 hex digits
	binary.BigEndian.PutUint16(response[8:10], code)
	binary.BigEndian.PutUint32(response[10:14], 2)           // number of apis
	binary.BigEndian.PutUint16(response[14:16], api_version) //aoi version
	response = append(response, 0x00, 0x00)                  // MinVersion = 0
	response = append(response, 0x00, 0x04)                  // MaxVersion = 4
	response = append(response, 0x00)                        // tagged fields
	response = append(response, 0x00, 0x00, 0x00, 0x00)      // throttle time
	response = append(response, 0x00)                        // tagged fields
	message_size = uint32(len(response))                     // The actual length of the response buffer excluding the first 4 bytes
	binary.BigEndian.PutUint32(response[0:4], message_size)  // Update message_size with the correct total length
	// response:= []byte{0,0,0,0,0,0,0,7} // just a hard coded way to send correlation id
	fmt.Println(response)

	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)

	}
	// fmt.Println("Here ", l)

}
