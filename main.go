package main

import (
	"context"
	"log"
	"net/http"

	"ogframe/pkg/socket"
	"ogframe/pkg/viewer"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	server := socket.NewServer("./images")

	ctx := context.Background()
	go server.Run(ctx)

	mux.Handle("/", viewer.NewViewer("frame"))
	mux.Handle("/socket", server)
	return http.ListenAndServe(":8080", mux)
}
