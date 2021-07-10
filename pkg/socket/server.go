package socket

import (
	"context"
	"encoding/base64"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

var _ http.Handler = (*Server)(nil)

type Server struct {
	join     chan *Client
	leave    chan *Client
	msgs     chan []byte
	clients  map[*Client]struct{}
	upgrader *websocket.Upgrader
	imgPath  string
}

func NewServer(imgPath string) *Server {
	return &Server{
		join:     make(chan *Client),
		leave:    make(chan *Client),
		msgs:     make(chan []byte),
		clients:  make(map[*Client]struct{}),
		upgrader: &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 256},
		imgPath:  imgPath,
	}
}

func (s *Server) Run(ctx context.Context) {
	go func() {
		for range time.Tick(4 * time.Second) {
			randFile, err := getRandomFileInDir(s.imgPath)
			if err != nil {
				// TODO: handle errors properly
				log.Print(err)
			}

			imgBytes, err := os.ReadFile(randFile)
			if err != nil {
				log.Print(err)
			}

			str := base64.StdEncoding.EncodeToString(imgBytes)
			s.msgs <- []byte(str)
		}
	}()
	for {
		select {
		case client := <-s.join:
			s.clients[client] = struct{}{}
		case client := <-s.leave:
			delete(s.clients, client)
			// TODO: is this the correct place to do this?
			close(client.send)
			close(client.receive)
		case msg := <-s.msgs:
			for client := range s.clients {
				client.send <- msg
			}
		case <-ctx.Done():
			return
		}
	}
}

// getRandomFileInDir returns the full path to a random file in the dir.
// https://stackoverflow.com/questions/45941821/how-do-you-get-full-paths-of-all-files-in-a-directory-in-go
func getRandomFileInDir(dir string) (string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	entry := dirEntries[rand.Intn(len(dirEntries))]
	fullEntryPath := filepath.Join(path, entry.Name())

	if entry.IsDir() {
		return getRandomFileInDir(fullEntryPath)
	}

	return fullEntryPath, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "upgrading websocket failed", http.StatusInternalServerError)
		return
	}

	client := &Client{
		socket:  socket,
		send:    make(chan []byte),
		receive: make(chan []byte),
	}
	s.join <- client
	defer func() { s.leave <- client }()
	go client.write()
	client.read()
}
