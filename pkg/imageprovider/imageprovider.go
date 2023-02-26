package imageprovider

import (
	"context"
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

//go:embed waiting.png
var waitingPicture []byte

type ImageProvider struct {
	dir string
	// files is a flat slice of all filepaths
	files map[string]struct{}

	m       *sync.RWMutex
	watcher *fsnotify.Watcher
}

func Run(ctx context.Context, dir string) (*ImageProvider, error) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ip := &ImageProvider{
		dir:     dir,
		watcher: watcher,
		m:       &sync.RWMutex{},
	}

	go func() {
		for {
			if err := ip.resetFilesOnChange(ctx); err != nil {
				return
			}
		}
	}()

	return ip, nil
}

func (ip *ImageProvider) Close() error {
	if ip == nil {
		return nil
	}

	return ip.watcher.Close()
}

func (ip *ImageProvider) Files() map[string]struct{} {
	ip.m.RLock()
	defer ip.m.RUnlock()
	return ip.files
}

func (ip *ImageProvider) reload() error {
	files, err := getFilePathsInDir(ip.dir)
	if err != nil {
		return err
	}

	ip.m.Lock()
	defer ip.m.Unlock()
	ip.files = files
	return nil
}

func (ip *ImageProvider) resetFilesOnChange(ctx context.Context) error {
	err := ip.watcher.Add(ip.dir)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-ip.watcher.Events:
		log.Println("files changed")
	case err := <-ip.watcher.Errors:
		log.Println("watched error: ", err)
	}

	ip.m.Lock()
	defer ip.m.Unlock()
	ip.files = nil
	return nil
}

func (ip *ImageProvider) GetNext() ([]byte, error) {
	if len(ip.Files()) == 0 {
		ip.reload()
	}

	var randomFile string
	for file := range ip.Files() {
		randomFile = file
		break
	}

	if randomFile == "" {
		return waitingPicture, nil
	}

	ip.m.Lock()
	defer ip.m.Unlock()
	delete(ip.files, randomFile)

	return os.ReadFile(randomFile)
}

// getFilePathsInDir returns a flat list of all files in a dir.
// It will go through the directories recursively and append just the contained files' paths.
func getFilePathsInDir(dir string) (map[string]struct{}, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := map[string]struct{}{}
	for _, entry := range dirEntries {
		fullEntryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			filesInDir, err := getFilePathsInDir(fullEntryPath)
			if err != nil {
				return nil, err
			}

			mergeMap(files, filesInDir)
			continue
		}

		files[fullEntryPath] = struct{}{}
	}

	return files, nil
}

func mergeMap(target, source map[string]struct{}) {
	for k, v := range source {
		target[k] = v
	}
}
