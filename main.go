package main

import (
	"go-sse/controller/event"
	r "go-sse/redis"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

const PORT = "8080"

func main() {

	if err := godotenv.Load(); err != nil {
		panic("failed to initilize dot env")
	}

	r.CreateRedisClient()

	mux := http.NewServeMux()

	event.RegisterEventHandler("/event", mux)

	staticDir := http.FileServer(http.Dir("./static"))
	mux.Handle("/", staticDir)

	log.Print("listening on port: ", PORT)
	err := http.ListenAndServe(":"+PORT, mux)
	if err != nil {
		log.Fatal(err)
	}

}
