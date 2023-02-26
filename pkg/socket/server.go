package socket

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/brumhard/pix/pkg/imageprovider"

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

func NewServer(imgPath string) *Server {
	return &Server{
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 256,
			// can be enabled if backend and frontend should run in different hosts
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		imgPath: imgPath,
	}
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
		// err will be promoted to client
		log.Print(err)
		return
	}

	defer socket.Close()

	s.sendImageLoop(r.Context(), socket, delay)
}

func (s *Server) sendImageLoop(ctx context.Context, socket *websocket.Conn, delay int) {
	images, err := imageprovider.Run(ctx, s.imgPath)
	if err != nil {
		log.Println(err)
	}

	for {
		imgBytes, err := images.GetNext()
		if err != nil {
			log.Print(err)
		}

		if err := socket.WriteMessage(websocket.BinaryMessage, imgBytes); err != nil {
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
