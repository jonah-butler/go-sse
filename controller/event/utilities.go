package event

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

func setSSEHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
}

func simulateLiveData(id string, ctx context.Context) {
	loopCeil := 10

	for i := 0; i < loopCeil; i++ {
		time.Sleep(1 * time.Second)
		data := id + "-" + strconv.Itoa(i)
		Broadcaster.Publish(ctx, id, data)
	}
}
