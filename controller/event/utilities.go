package event

import (
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

func simulateLiveData(eventCh *chan string, id string) {
	loopCeil := 10

	for i := 0; i < loopCeil; i++ {
		time.Sleep(1 * time.Second)
		*eventCh <- id + "-" + strconv.Itoa(i)
	}
}
