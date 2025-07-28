package core

import (
	"errors"
	"os"
	"sync"
)

type entry struct {
	K string  `json:"k"`
	V *string `json:"V"`
}

type DB struct {
	mu sync.RWMutex
	// key -> file offset
	index map[string]int64
	f     *os.File
	path  string
}

var (
	ErrKeyNotFound = errors.New("key not found")
)
