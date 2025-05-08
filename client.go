package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Failed to connect:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server!")

	// Read welcome message
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println("Server:", scanner.Text())
		}
	}()

	// Send user input to server
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("You: ")
		if input.Scan() {
			text := input.Text()
			if text == "exit" {
				break
			}
			conn.Write([]byte(text + "\n"))
		}
	}
}
