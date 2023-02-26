package main

import (
	"flag"
	"io/fs"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"ogframe/frontend"
	ownhttp "ogframe/pkg/http"
	"ogframe/pkg/socket"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var imgPath string
	flag.StringVar(&imgPath, "images", "", "path to images to be shown on the frame")
	flag.Parse()

	if imgPath == "" {
		return errors.New("images argument needs to be set")
	}

	router := http.NewServeMux()
	dist, err := fs.Sub(frontend.Static, "build/web")
	if err != nil {
		return err
	}

	router.Handle("/api/socket", socket.NewServer(imgPath))
	router.Handle("/", ownhttp.NewSPAHandler(http.FS(dist), "index.html"))

	log.Print("server starting")

	return http.ListenAndServe(":8080", router)
}
