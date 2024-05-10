package sse

import (
	"context"
	"fmt"
	r "go-sse/redis"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	ID      string
	Channel chan string
}

type Broadcaster struct {
	M  sync.Mutex
	S  map[string][]Subscriber
	PS map[string]*redis.PubSub
}

// generates a single broadcaster to manage
// connections during application lifetime
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		S:  make(map[string][]Subscriber),
		PS: make(map[string]*redis.PubSub),
	}
}

func (b *Broadcaster) Subscribe(ctx context.Context, userID string) Subscriber {
	subscriber := Subscriber{
		ID:      userID,
		Channel: make(chan string), // unbuffered for now
	}

	b.M.Lock()
	b.S[userID] = append(b.S[userID], subscriber)
	fmt.Println(b.S[userID])

	if _, ok := b.PS[userID]; !ok {
		b.PS[userID] = r.CreatePubSubClient(ctx, userID)
		go b.Broadcast(subscriber) // only sub a single user connection
	}

	defer b.M.Unlock()
	return subscriber
}

func (b *Broadcaster) Unsubscribe(ctx context.Context, userID string, subscriber Subscriber) {
	b.M.Lock()
	defer b.M.Unlock()
	subsuscribers := b.S[userID]
	pubsub := b.PS[userID]

	for i, s := range subsuscribers {
		if s.Channel == subscriber.Channel {
			b.S[userID] = append(subsuscribers[:i], subsuscribers[i+1:]...)
			break
		}
	}

	close(subscriber.Channel)

	if len(b.S[userID]) == 0 {
		delete(b.S, userID)
		pubsub.Unsubscribe(ctx)
		pubsub.Close()
		delete(b.PS, userID)
	}
}

func (b *Broadcaster) Publish(ctx context.Context, userID, data string) {
	b.M.Lock()
	defer b.M.Unlock()

	for range b.S[userID] {
		r.PublishToRedis(ctx, userID, data)
	}
}

func (b *Broadcaster) Broadcast(subscriber Subscriber) {
	b.M.Lock()
	pubsub := b.PS[subscriber.ID]
	b.M.Unlock()

	for message := range pubsub.Channel() {
		b.M.Lock()
		for i, sub := range b.S[subscriber.ID] {
			fmt.Println(i)
			sub.Channel <- message.Payload
		}
		b.M.Unlock()
	}
}
