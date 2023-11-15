package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	writer := bufio.NewWriter(conn)

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error:", err)
			return
		}

		fmt.Printf("Received: %s\n", buffer[:n])
		_, err = writer.WriteString(fmt.Sprintf("You send: %s", string(buffer)))
		if err != nil {
			fmt.Println("Error Write:", err)
			return
		}

		if err := writer.Flush(); err != nil {
			fmt.Println("Error Flush:", err)
			return
		}
	}

}
