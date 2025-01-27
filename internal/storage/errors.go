package storage

import "errors"

var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrDataCorruption  = errors.New("data corruption detected")
	ErrIndexClosed     = errors.New("index manager is closed")
	ErrInvalidArgument = errors.New("invalid arguments")
)
