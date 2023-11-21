package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/samuelsih/chat-encrypt/des"
	"github.com/samuelsih/chat-encrypt/server/worker"
)

var clients = worker.NewClients()

func exec(conn net.Conn, msg string) error {
	input := sanitize(msg)
	cmd := strings.Split(input, " ")
	remoteAddr := conn.RemoteAddr().String()

	switch cmd[0] {
	case "USERNAME":
		username, exist := clients.FindUsernameByRemoteAddr(remoteAddr)
		if !exist {
			clients.AddClient(cmd[1], conn)
			return clients.Broadcast(fmt.Sprintf("RESPUSR User %s has joined\n", cmd[1]))
		}

		return singleSend(conn, fmt.Sprintf("ERROR Already registered as %s\n", username))

	case "MESSAGE":
		username, exist := clients.FindUsernameByRemoteAddr(remoteAddr)
		if !exist {
			return singleSend(conn, "ERROR Unregistered\n")
		}

		encryptedMsg := des.Encrypt(strings.Join(cmd[1:], " "), des.EncryptionBase64)
		return clients.Broadcast(fmt.Sprintf("RESPMSG %s: %s\n", username, encryptedMsg))

	case "EXIT":
		if len(cmd) > 1 {
			return singleSend(conn, "Unknown command. Please write correct command\n")
		}

		_, exist := clients.FindUsernameByRemoteAddr(remoteAddr)
		if !exist {
			return singleSend(conn, "ERROR Unknown user\n")
		}

		return singleSend(conn, "EXIT\n")
	default:
		return singleSend(conn, "Unknown command. Please write correct command\n")
	}
}

func sanitize(s string) string {
	trimSpaced := strings.TrimSpace(s)
	trimRight := strings.TrimRight(trimSpaced, "\r\n")
	return strings.TrimRight(trimRight, "\n")
}
