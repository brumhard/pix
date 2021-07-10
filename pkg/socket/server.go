package socket

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"ogframe/pkg/fileindex"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	delayKey     = "delay"
	defaultDelay = 5
)

var _ http.Handler = (*Server)(nil)

type Server struct {
	upgrader *websocket.Upgrader
	imgPath  string
}

func NewServer(imgPath string) (*Server, error) {

	return &Server{
		upgrader: &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 256},
		imgPath:  imgPath,
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL)

	delay := defaultDelay
	delayStr := r.URL.Query().Get(delayKey)
	if delayStr != "" {
		var err error
		delay, err = strconv.Atoi(delayStr)
		if err != nil {
			http.Error(w, "expected int type as delay", http.StatusBadRequest)
			return
		}
	}

	socket, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "upgrading websocket failed", http.StatusInternalServerError)
		return
	}

	defer socket.Close()

	s.sendImageLoop(r.Context(), socket, delay)
}

func (s *Server) sendImageLoop(ctx context.Context, socket *websocket.Conn, delay int) {
	fi, err := fileindex.New(s.imgPath)
	if err != nil {
		log.Print(err)
	}

	for {
		randFile, err := fi.GetRandomFile()
		if err != nil {
			log.Print(err)
		}

		imgBytes, err := os.ReadFile(randFile)
		if err != nil {
			log.Print(err)
		}

		str := base64.StdEncoding.EncodeToString(imgBytes)

		if err := socket.WriteMessage(websocket.BinaryMessage, []byte(str)); err != nil {
			if websocket.IsCloseError(
				err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived,
			) {
				return
			}
			log.Print(err)
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(delay) * time.Second):
		}
	}
}
