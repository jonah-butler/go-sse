package event

import (
	"fmt"
	"go-sse/sse"
	"net/http"
)

var Broadcaster = sse.NewBroadcaster()

func handleEventIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("userId")

	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	setSSEHeaders(w)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Failed to write data", http.StatusExpectationFailed)
		return
	}

	subscriber := Broadcaster.Subscribe(ctx, userID)
	defer Broadcaster.Unsubscribe(ctx, userID, subscriber)

	// simulate data being sent over a subscribed channel.
	go simulateLiveData(userID, ctx)

	for {
		select {

		case <-ctx.Done():
			fmt.Println("context closed")
			return

		case data := <-subscriber.Channel:
			fmt.Println("got data...")
			_, err := fmt.Fprintf(w, "data: %s\n\n", data)
			if err != nil {
				http.Error(w, "failed to write data to channel(s): "+err.Error(), http.StatusBadRequest)
				return
			}
			flusher.Flush()
		}
	}

}
