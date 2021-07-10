package fileindex

import (
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

type FileIndex struct {
	// files is a flat slice of all filepaths
	files []string
	// usedIndeces keeps track of what files have been already selected
	unusedIndeces []int
	m             *sync.RWMutex
}

func New(dir string) (*FileIndex, error) {
	// files is a flattened list of all files in the directory.
	// This is done to ensure that all the images have the same chance to be picked when using a random index.
	files, err := getFilePathsInDir(dir)
	if err != nil {
		return nil, err
	}

	return &FileIndex{
		files:         files,
		unusedIndeces: initIndexSlice(len(files)),
		m:             &sync.RWMutex{},
	}, nil
}

// getFilePathsInDir returns a flat list of all files in a dir.
// It will go through the directories recursively and append just the contained files' paths.
func getFilePathsInDir(dir string) ([]string, error) {
	path, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range dirEntries {
		fullEntryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			filesInDir, err := getFilePathsInDir(fullEntryPath)
			if err != nil {
				return nil, err
			}

			files = append(files, filesInDir...)
			continue
		}
		files = append(files, fullEntryPath)
	}

	return files, nil
}

func (fi *FileIndex) GetRandomFile() (string, error) {
	fi.m.Lock()
	defer fi.m.Unlock()

	// if all indixes have been used reinitialize the slice and start over
	if len(fi.unusedIndeces) == 0 {
		fi.unusedIndeces = initIndexSlice(len(fi.files))
	}

	// get a random index for the index slice
	unusedI := rand.Intn(len(fi.unusedIndeces))
	i := fi.unusedIndeces[unusedI]
	// remove the index from the slice after it has been used
	fi.unusedIndeces = removeFromIntSliceByIndex(fi.unusedIndeces, unusedI)

	return fi.files[i], nil
}

// removeFromIntSliceByIndex removes the item at index i from the slice s.
// The order of the slice will not be maintained.
func removeFromIntSliceByIndex(s []int, i int) []int {
	// overwrite the item at index i with the last item in the slice
	s[i] = s[len(s)-1]
	// return subslice with last item omitted
	return s[:len(s)-1]
}

func initIndexSlice(len int) []int {
	slice := make([]int, 0, len)
	for i := 0; i < len; i++ {
		slice = append(slice, i)
	}

	return slice
}
