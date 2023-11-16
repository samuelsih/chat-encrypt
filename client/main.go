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

	go prompt(conn)
	go recv(conn)

	wg.Wait()
}

func prompt(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if ok := scanner.Scan(); !ok {
			if err := scanner.Err(); err != nil {
				fmt.Println("Error scanner:", err)
				return
			}
		}

		msg := fmt.Sprintf("%s\n", scanner.Text())

		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error conn.Write:", err)
			return
		}
	}
}

func recv(listener net.Conn) {
	defer wg.Done()
	reader := bufio.NewReader(listener)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("[x]", msg)
	}
}
