package sse

import (
	"context"
	r "go-sse/redis"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
)

type ConnectionDetails struct {
	Channel     chan string
	PubSub      *redis.PubSub
	Connections int
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

	connectionDetails, ok := conn.CD[id]
	if !ok {
		c = make(chan string)

		cd := ConnectionDetails{
			Channel: c,
			PubSub:  r.CreatePubSubClient(ctx, id),
		}

		connectionDetails = cd
	}

	connectionDetails.Connections++
	conn.CD[id] = connectionDetails

	return &c
}

func (conn *SSEConn) RemoveClient(ctx context.Context, id string) {
	conn.M.Lock()
	defer conn.M.Unlock()

	cd, ok := conn.CD[id]
	if !ok {
		return
	}

	cd.Connections -= 1

	if cd.Connections == 0 {
		if err := cd.PubSub.Unsubscribe(ctx); err != nil {
			log.Fatal("failed to unsubscribe from pub sub channel")
		}
		if err := cd.PubSub.Close(); err != nil {
			log.Fatal("failed to close pub sub channel")
		}

		close(cd.Channel)
		delete(conn.CD, id)
	}
}

func (conn *SSEConn) ListenOnChannel(id string) {
	client := conn.CD[id]
	for msg := range client.PubSub.Channel() {
		client.Channel <- msg.Payload
	}
}
