package storage

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	IndexFileName = "lumora.index"
)

type IndexEntry struct {
	Offset    int64
	KeySize   uint32
	ValueSize uint32
	Timestamp uint32
}

type IndexManager struct {
	mu       sync.RWMutex
	index    map[string]IndexEntry
	file     *os.File
	filepath string
	closed   bool
}

func NewIndexManager(path string) (*IndexManager, error) {

	filepath := fmt.Sprintf("%s/%s", path, IndexFileName)

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to load index file: %w", err)
	}

	im := &IndexManager{
		index:    make(map[string]IndexEntry),
		file:     file,
		filepath: filepath,
	}

	if err := im.load(); err != nil {
		return nil, fmt.Errorf("failed to load index: %w", err)
	}
	return im, nil
}

func (im *IndexManager) Save() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.closed {
		return ErrIndexClosed
	}

	if err := im.file.Truncate(0); err != nil {
		return fmt.Errorf("truncate failed: %w", err)
	}

	if _, err := im.file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek failed: %w", err)
	}

	if err := gob.NewEncoder(im.file).Encode(im.index); err != nil {
		return fmt.Errorf("gob encode failed: %w", err)
	}

	return im.file.Sync()
}

func (im *IndexManager) load() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if _, err := im.file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if err := gob.NewDecoder(im.file).Decode(&im.index); err != nil {
		if err == io.EOF {
			return nil
		}
		return fmt.Errorf("gob decode failed: %w", err)
	}
	return nil
}

func (im *IndexManager) GetEntry(key string) (IndexEntry, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	if im.closed {
		return IndexEntry{}, ErrIndexClosed
	}

	entry, exists := im.index[key]
	if !exists {
		return IndexEntry{}, ErrKeyNotFound
	}

	return entry, nil
}

func (im *IndexManager) PutEntry(key string, entry IndexEntry) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.closed {
		return ErrIndexClosed
	}

	im.index[key] = entry
	return nil
}

func (im *IndexManager) DeleteEntry(key string) error {
  im.mu.Lock()
  defer im.mu.Unlock()

  if im.closed{
    return ErrIndexClosed
  }

  delete(im.index, key)
  return nil
}

func (im *IndexManager) Close() error {
	im.mu.Lock()
	defer im.mu.Unlock()

	if im.closed {
		return nil
	}

	if err := im.file.Close(); err != nil {
		return fmt.Errorf("index file close failed: %w", err)
	}

	im.closed = true
	return nil
}
