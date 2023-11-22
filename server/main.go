package main

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/samuelsih/chat-encrypt/des"
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

		fmt.Println("Connection established", conn.RemoteAddr().String())
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error:", err)
				return
			}
		}

		if len(msg) <= 1 {
			return
		}

		decryptedTxt := des.Decrypt(msg, des.DecryptionBase64)
		exec(conn, decryptedTxt)
	}
}
