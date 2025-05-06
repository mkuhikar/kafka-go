package main

import (
	"bytes"
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
	fmt.Printf("the no of bytes received from client are %d \n", n) /*fmt.Printf allows formatted strings, just like printf in C.


	%d is a placeholder for an integer — here, it's replaced with the value of n.

	\n adds a newline after the output. */
	fmt.Println("Raw bytes in hex:")
	for i := 0; i < n; i++ {
		fmt.Printf("%02X ", buffer[i])
		if (i+1)%16 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()

	if len(buffer) < 12 {
		fmt.Println("Request too short to contain a valid correlation ID")
		os.Exit(1)
	}
	// correlation_id := binary.BigEndian.Uint32(buffer[4:8]) //converted 4 bytes of buffer to uInt32, in big endian most significant byte comes first
	apiKey := binary.BigEndian.Uint16(buffer[4:6])
	apiVersion := binary.BigEndian.Uint16(buffer[6:8])
	correlation_id := binary.BigEndian.Uint32(buffer[4:8])

	fmt.Printf("API Key: %d\n", apiKey)
	fmt.Printf("API Version: %d\n", apiVersion)
	fmt.Printf("Correlation ID: %d\n", correlation_id)
	fmt.Printf("correlation id %d \n", correlation_id)
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
	var response bytes.Buffer

	// Correlation ID
	binary.Write(&response, binary.BigEndian, int32(correlation_id))

	// Error Code (set to 0)
	binary.Write(&response, binary.BigEndian, int16(0))

	// Number of API keys (set to 2, should be INT8)
	binary.Write(&response, binary.BigEndian, int8(2))

	// API Key Entry for ApiVersions (18)
	binary.Write(&response, binary.BigEndian, int16(18)) // api_key
	binary.Write(&response, binary.BigEndian, int16(0))  // min_version
	binary.Write(&response, binary.BigEndian, int16(4))  // max_version

	// Tagged fields (set to INT8, always 0)
	binary.Write(&response, binary.BigEndian, int8(0))

	// Throttle time (set to 0)
	binary.Write(&response, binary.BigEndian, int32(0))

	// Tagged fields (set to INT8, always 0) after throttle_time_ms
	binary.Write(&response, binary.BigEndian, int8(0))

	// Now fix message size (total - 4 bytes for the size itself)
	messageSize := make([]byte, 4)
	binary.BigEndian.PutUint32(messageSize, uint32(response.Len())) // Include the message size
	if _, err := conn.Write(messageSize); err != nil {
		fmt.Println("Error writing response size: ", err.Error())
		os.Exit(1)
	}

	// Write the full response
	if _, err := conn.Write(response.Bytes()); err != nil {
		fmt.Println("Error writing response: ", err.Error())
		os.Exit(1)
	}
	// fmt.Println("Here ", l)

}
