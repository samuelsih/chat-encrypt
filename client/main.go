package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(1)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer conn.Close()

	go recv(conn)
	go prompt(conn)

	wg.Wait()
}

func prompt(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				fmt.Println("Error scanner:", err)
				wg.Done()
				return
			}
		}

		_, err := conn.Write([]byte(scanner.Text()))
		if err != nil {
			fmt.Println("Error conn.Write:", err)
			wg.Done()
			return
		}
	}
}

func recv(listener net.Conn) {
	for {
		buffer := make([]byte, 1024)
		_, err := listener.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println("Error:", err)
			break
		}

		fmt.Println(string(buffer))
	}
}
