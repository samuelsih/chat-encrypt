package worker

import (
	"fmt"
	"net"
	"time"
)

func (m *clients[K, V]) detectDisconnect(username K, conn V) {
	disconnect := make(chan struct{})
	ticker := time.NewTicker(time.Second * 1)

	defer close(disconnect)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ping(conn, disconnect)
		case <-disconnect:
			m.deleteByUsername(username)
			m.Broadcast(fmt.Sprintf("LEAVE User %s has leave from server\n", username))
			return
		}
	}
}

func ping(conn net.Conn, disconnect chan<- struct{}) {
	go func() {
		n, err := conn.Write([]byte("PING\n"))
		if err != nil || n == 0 {
			disconnect <- struct{}{}
			conn.Close()
		}
	}()
}

func (m *clients[K, _]) deleteByUsername(username K) {
	m.internal.Delete(username)
}
