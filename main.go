package main

import (
	"crypto/subtle"
	"errors"
	"flag"
	"fmt"
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
	var creds string
	flag.StringVar(&imgPath, "images", "", "path to images to be shown on the frame")
	flag.StringVar(&creds, "credentials", "", "comma delimited user:password for webdav server access")
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
	router.Handle(webdavServer.Prefix, basicAuthMW(webdavServer, creds))
	router.Handle("/", ownhttp.NewSPAHandler(http.FS(dist), "index.html"))

	log.Print("server starting")
	go http.ListenAndServe(":8420", router)

	signalc := make(chan os.Signal, 1)
	signal.Notify(signalc, syscall.SIGINT, syscall.SIGTERM)
	<-signalc

	return nil
}

func basicAuthMW(next http.Handler, credentials string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if credentials != "" {
			user, pw, ok := r.BasicAuth()
			providedCreds := fmt.Sprintf("%s:%s", user, pw)
			if !ok || subtle.ConstantTimeCompare([]byte(providedCreds), []byte(credentials)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="get good in"`)
				w.WriteHeader(401)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
