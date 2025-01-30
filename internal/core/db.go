package core

import (
	"fmt"
	"sync"
	"time"

	"github.com/CinematicCow/lumora/internal/storage"
)

type LumoraDB struct {
	dataManager  *storage.DataManager
	indexManager *storage.IndexManager
	mu           sync.RWMutex
}

func Open(dataDir string) (*LumoraDB, error) {
	dataManager, err := storage.NewDataManager(dataDir)
	if err != nil {
		return nil, fmt.Errorf("data manager init failed: %w", err)
	}

	indexManager, err := storage.NewIndexManager(dataDir)
	if err != nil {
		dataManager.Close()
		return nil, fmt.Errorf("index manager init failed: %w", err)
	}

	return &LumoraDB{
		dataManager:  dataManager,
		indexManager: indexManager,
	}, nil
}

func (db *LumoraDB) Put(key string, value []byte) error {
	if key == "" {
		return storage.ErrInvalidArgument
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	record := &storage.DataRecord{
		Timestamp: uint32(time.Now().Unix()),
		KeySize:   uint32(len(key)),
		ValueSize: uint32(len(value)),
		Key:       []byte(key),
		Value:     value,
	}

	offset, err := db.dataManager.WriteRecord(record)
	if err != nil {
		return fmt.Errorf("write record failed: %w", err)
	}

	entry := storage.IndexEntry{
		Offset:    offset,
		KeySize:   record.KeySize,
		ValueSize: record.ValueSize,
		Timestamp: record.Timestamp,
	}

	if err := db.indexManager.PutEntry(key, entry); err != nil {
		return fmt.Errorf("index update failed: %w", err)
	}

	if err := db.indexManager.Save(); err != nil {
		return fmt.Errorf("index save failed: %w", err)
	}

	return nil
}

func (db *LumoraDB) Get(key string) ([]byte, error) {
	if key == "" {
		return nil, storage.ErrInvalidArgument
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	entry, err := db.indexManager.GetEntry(key)
	if err != nil {
		return nil, err
	}

	record, err := db.dataManager.ReadRecord(entry.Offset)
	if err != nil {
		return nil, fmt.Errorf("data read failed: %w", err)
	}

	// validate record integrity
	if string(record.Key) != key {
		return nil, fmt.Errorf("%w: key mismatch", storage.ErrDataCorruption)
	}

	if uint32(len(record.Value)) != entry.ValueSize {
		return nil, fmt.Errorf("%w: value size mismatch", storage.ErrDataCorruption)
	}
	return record.Value, nil
}

func (db *LumoraDB) Close() error {

	var errs []error

	if err := db.dataManager.Close(); err != nil {
		errs = append(errs, fmt.Errorf("data manager close: %w", err))
	}

	if err := db.indexManager.Close(); err != nil {
		errs = append(errs, fmt.Errorf("index manager close: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("database close errors: %v", errs)
	}
	return nil
}
