package sse

import (
	"context"
	r "go-sse/redis"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type ConnectionDetails struct {
	C  chan string
	PS *redis.PubSub
}

type SSEConn struct {
	M  sync.Mutex
	CD map[string]ConnectionDetails
}

func NewSSECon() *SSEConn {
	return &SSEConn{CD: make(map[string]ConnectionDetails)}
}

func (conn *SSEConn) AddClient(ctx context.Context, id string) *chan string {
	var c chan string

	conn.M.Lock()
	defer conn.M.Unlock()

	c = conn.CD[id].C
	if c == nil {
		c = make(chan string)

		cd := ConnectionDetails{
			C:  c,
			PS: r.CreatePubSubClient(ctx, id),
		}

		conn.CD[id] = cd
	}

	return &c
}

func (conn *SSEConn) RemoveClient(ctx context.Context, id string) {
	conn.M.Lock()
	defer conn.M.Unlock()

	cd, ok := conn.CD[id]
	if !ok {
		return
	}

	if err := cd.PS.Unsubscribe(ctx); err != nil {
		log.Fatal("failed to unsubscribe from pub sub channel")
	}
	if err := cd.PS.Close(); err != nil {
		log.Fatal("failed to close pub sub channel")
	}

	close(cd.C)
	delete(conn.CD, id)
}

func (conn *SSEConn) ListenOnChannel(id string) {
	client := conn.CD[id]
	for msg := range client.PS.Channel() {
		client.C <- msg.Payload
	}
}
