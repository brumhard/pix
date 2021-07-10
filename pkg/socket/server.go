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
	upgrader *websocket.Upgrader
	imgPath  string
}

func NewServer(imgPath string) *Server {
	return &Server{
		upgrader: &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 256},
		imgPath:  imgPath,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Print(r.URL)
	socket, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "upgrading websocket failed", http.StatusInternalServerError)
		return
	}

	s.sendImageLoop(r.Context(), socket, 4)
}

func (s *Server) sendImageLoop(ctx context.Context, socket *websocket.Conn, delay int) {
	for {
		randFile, err := getRandomFileInDir(s.imgPath)
		if err != nil {
			log.Print(err)
		}

		imgBytes, err := os.ReadFile(randFile)
		if err != nil {
			log.Print(err)
		}

		str := base64.StdEncoding.EncodeToString(imgBytes)

		if err := socket.WriteMessage(websocket.BinaryMessage, []byte(str)); err != nil {
			log.Print(err)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(delay) * time.Second):
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
