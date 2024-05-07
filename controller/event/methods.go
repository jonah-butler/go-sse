package event

import (
	"encoding/json"
	"fmt"
	"go-sse/sse"
	"net/http"
)

var CONNECTIONS = sse.NewSSECon()

type User struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

func handleEventIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := r.URL.Query().Get("userId")
	user := &User{
		Id: userId,
	}

	if userId == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	setSSEHeaders(w)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Failed to write data", http.StatusExpectationFailed)
		return
	}

	ch := CONNECTIONS.AddClient(ctx, userId)
	go CONNECTIONS.ListenOnChannel(userId)

	defer CONNECTIONS.RemoveClient(ctx, userId)

	// simulate data being sent over a subscribed channel.
	go simulateLiveData(userId, ctx)

	for {
		select {
		// listen for closed connections
		case <-ctx.Done():
			return
		// listen for data on user subscribed channel
		case data := <-*ch:
			user.Data = data
			d, err := json.Marshal(user)
			if err != nil {
				http.Error(w, "Failed to marshal json", http.StatusBadRequest)
				return
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", string(d))
			if err != nil {
				http.Error(w, "Failed to write data", http.StatusBadRequest)
				return
			}
			flusher.Flush()
		}
	}

}
