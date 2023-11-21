package worker

import (
	"bufio"
	"net"
	"sync"
)

func NewClients() clients[string, net.Conn] {
	return clients[string, net.Conn]{}
}

type clients[K string, V net.Conn] struct {
	internal sync.Map
}

func (m *clients[K, V]) AddClient(username K, conn V) {
	m.internal.Store(username, conn)
	go m.detectDisconnect(username, conn)
}

func (m *clients[K, _]) UsernameExist(username K) bool {
	_, exists := m.internal.Load(username)
	return exists
}

func (m *clients[_, _]) FindUsernameByRemoteAddr(addr string) (string, bool) {
	var (
		exists   bool
		username string
	)

	m.internal.Range(func(key, value any) bool {
		conn := value.(net.Conn)
		if conn.RemoteAddr().String() == addr {
			exists = true
			username = key.(string)
			return false
		}

		return true
	})

	return username, exists
}

func (m *clients[_, _]) Broadcast(msg string) error {
	var broadcastErr error

	m.internal.Range(func(key, value any) bool {
		conn := value.(net.Conn)

		writer := bufio.NewWriter(conn)
		_, err := writer.WriteString(msg)
		if err != nil {
			broadcastErr = err
			return false
		}

		if err := writer.Flush(); err != nil {
			broadcastErr = err
			return false
		}

		return true
	})

	return broadcastErr
}
