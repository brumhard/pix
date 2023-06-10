package main

import (
	"errors"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/brumhard/pix/frontend"
	ownhttp "github.com/brumhard/pix/pkg/http"
	"github.com/brumhard/pix/pkg/socket"
	ownwebdav "github.com/brumhard/pix/pkg/webdav"
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

	webdavServer, err := ownwebdav.NewServer(imgPath, "/store/")
	if err != nil {
		return err
	}
	defer webdavServer.Close()

	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.Handle("/api/socket", socket.NewServer(imgPath))
	router.Handle(webdavServer.Prefix, webdavServer)
	router.Handle("/", ownhttp.NewSPAHandler(http.FS(dist), "index.html"))

	log.Print("server starting")
	go http.ListenAndServe(":8420", router)

	signalc := make(chan os.Signal, 1)
	signal.Notify(signalc, syscall.SIGINT, syscall.SIGTERM)
	<-signalc

	return nil
}
