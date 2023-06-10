package local

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/net/webdav"
)

var _ webdav.FileSystem = (*LocalFS)(nil)
var _ io.Closer = (*LocalFS)(nil)

type LocalFS struct {
	root string
}

func NewLocalFS(rootDir string) (*LocalFS, error) {
	if rootDir == "" {
		return nil, fmt.Errorf("rootDir is required")
	}

	root, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, err
	}

	return &LocalFS{
		root: root,
	}, nil
}

// Close implements io.Closer.
func (lfs *LocalFS) Close() error {
	return filepath.WalkDir(lfs.root, fs.WalkDirFunc(func(path string, d fs.DirEntry, _ error) error {
		if name := d.Name(); strings.HasPrefix(name, "._") {
			_ = os.RemoveAll(path)
		}
		return nil
	}))
}

// Mkdir implements webdav.FileSystem.
func (lfs *LocalFS) Mkdir(_ context.Context, name string, perm fs.FileMode) error {
	name = path.Join(lfs.root, name)
	return os.Mkdir(name, perm)
}

// OpenFile implements webdav.FileSystem.
func (lfs *LocalFS) OpenFile(_ context.Context, name string, flag int, perm fs.FileMode) (webdav.File, error) {
	name = path.Join(lfs.root, name)
	return os.OpenFile(name, flag, perm)
}

// RemoveAll implements webdav.FileSystem.
func (lfs *LocalFS) RemoveAll(_ context.Context, name string) error {
	name = path.Join(lfs.root, name)
	return os.RemoveAll(name)
}

// Rename implements webdav.FileSystem.
func (lfs *LocalFS) Rename(_ context.Context, oldName string, newName string) error {
	oldName = path.Join(lfs.root, oldName)
	newName = path.Join(lfs.root, newName)
	return os.Rename(oldName, newName)
}

// Stat implements webdav.FileSystem.
func (lfs *LocalFS) Stat(_ context.Context, name string) (fs.FileInfo, error) {
	name = path.Join(lfs.root, name)
	return os.Stat(name)
}
