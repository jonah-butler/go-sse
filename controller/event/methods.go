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
		fmt.Println("sse not supported")
		return
	}

	ch := CONNECTIONS.AddClient(userId)
	defer CONNECTIONS.RemoveClient(userId)

	go simulateLiveData(ch, userId)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("connection cancelled")
			return
		case data := <-*ch:
			user.Data = data
			d, err := json.Marshal(user)
			if err != nil {
				http.Error(w, "Failed to marshal json", http.StatusBadRequest)
				return
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", string(d))
			if err != nil {
				fmt.Println("err writing daata: ", err.Error())
				return
			}
			flusher.Flush()
		}
	}

}
