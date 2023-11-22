package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/kyokomi/emoji/v2"
	"github.com/samuelsih/chat-encrypt/des"
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

		msg := fmt.Sprintf("%s", scanner.Text())
		encryptedTxt := des.Encrypt(msg, des.EncryptionBase64)

		_, err := conn.Write([]byte(encryptedTxt + "\n"))
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
		response, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Error:", err)
			return
		}

		if response == "PING\n" {
			continue
		}

		if response == "EXIT\n" {
			return
		}

		msg := resp(response)

		fmt.Println(msg)
	}
}

func resp(msg string) string {
	input := sanitize(msg)
	cmd := strings.Split(input, " ")
	defaultMsg := strings.Join(cmd[1:], " ")

	switch cmd[0] {
	case "RESPUSR":
		return emoji.Sprintf(":man: %s", defaultMsg)
	case "ERROR":
		return emoji.Sprintf(":error: %s", defaultMsg)
	case "LEAVE":
		return emoji.Sprintf(":door: %s", defaultMsg)
	case "RESPMSG":
		decryptedMsg := des.Decrypt(cmd[2], des.DecryptionBase64)
		return emoji.Sprintf(":kiss: %s %s", cmd[1], decryptedMsg)
	default:
		return emoji.Sprint(":gorilla: Invalid message protocol sent by server")
	}
}

func sanitize(s string) string {
	trimSpaced := strings.TrimSpace(s)
	trimRight := strings.TrimRight(trimSpaced, "\r\n")
	return strings.TrimRight(trimRight, "\n")
}
