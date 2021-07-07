package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"ogframe/frontend"

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

	dist, err := fs.Sub(frontend.Static, "dist")
	if err != nil {
		return err
	}

	mux.Handle("/", http.FileServer(http.FS(dist)))
	mux.Handle("/old", viewer.NewViewer("frame"))
	mux.Handle("/socket", server)

	return http.ListenAndServe(":8080", mux)
}
