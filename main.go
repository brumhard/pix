package main

import (
	"context"
	"io/fs"
	"log"
	"net/http"
	"ogframe/frontend"
	ownhttp "ogframe/pkg/http"

	"github.com/gorilla/mux"

	"ogframe/pkg/socket"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	router := mux.NewRouter()
	server := socket.NewServer("./images")

	ctx := context.Background()
	go server.Run(ctx)

	dist, err := fs.Sub(frontend.Static, "dist")
	if err != nil {
		return err
	}

	router.Handle("/api/socket", server)
	router.PathPrefix("/").Handler(ownhttp.NewSPAHandler(http.FS(dist), "index.html"))

	return http.ListenAndServe(":8080", router)
}
