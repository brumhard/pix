package main

import (
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"

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
	imgPath := flag.String("images", "", "path to images to be shown on the frame")
	flag.Parse()

	if imgPath == nil || *imgPath == "" {
		return errors.New("images argument needs to be set")
	}

	_, err := os.Stat(*imgPath)
	if err != nil {
		return errors.Wrap(err, "images argument should point to a valid path")
	}

	router := http.NewServeMux()

	server := socket.NewServer(*imgPath)

	dist, err := fs.Sub(frontend.Static, "build/web")
	if err != nil {
		return err
	}

	router.Handle("/api/socket", server)
	router.Handle("/", ownhttp.NewSPAHandler(http.FS(dist), "index.html"))

	log.Print("server starting")

	return http.ListenAndServe(":8080", router)
}
