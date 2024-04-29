package event

import (
	"net/http"
)

func RegisterEventHandler(prefix string, server *http.ServeMux) {
	server.HandleFunc(prefix+"/sse", handleEventIndex)
}
