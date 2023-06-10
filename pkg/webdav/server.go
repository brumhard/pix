package local

import (
	"golang.org/x/net/webdav"
)

type Server struct {
	webdav.Handler
	localFS *LocalFS
}

func (s *Server) Close() error {
	if s.localFS == nil {
		return nil
	}
	return s.localFS.Close()
}

func NewServer(rootDir string, handlerPath string) (*Server, error) {
	fs, err := NewLocalFS(rootDir)
	if err != nil {
		return nil, err
	}

	return &Server{
		Handler: webdav.Handler{
			Prefix:     handlerPath,
			FileSystem: fs,
			LockSystem: webdav.NewMemLS(),
		},
		localFS: fs,
	}, nil
}
