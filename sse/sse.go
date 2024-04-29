package sse

import (
	"fmt"
	"sync"
)

type SSEConn struct {
	M sync.Mutex
	C map[string]chan string
}

func NewSSECon() *SSEConn {
	return &SSEConn{C: make(map[string]chan string)}
}

func (conn *SSEConn) AddClient(id string) *chan string {
	var c chan string

	conn.M.Lock()
	defer conn.M.Unlock()

	c, ok := conn.C[id]
	if !ok {
		c = make(chan string)
		conn.C[id] = c
	}

	return &c
}

func (conn *SSEConn) RemoveClient(id string) {
	fmt.Println("removing client..")
	conn.M.Lock()
	defer conn.M.Unlock()

	chann, ok := conn.C[id]
	if !ok {
		return
	}

	close(chann)
	delete(conn.C, id)
	fmt.Println("current list of clients..", conn.C)
}
