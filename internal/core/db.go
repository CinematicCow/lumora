package core

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Open(dir string) (*DB, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	path := filepath.Join(dir, "data.jsonl")

	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	db := &DB{
		index: make(map[string]int64),
		f:     f,
		path:  path,
	}

	if err := db.buildindex(); err != nil {
		closeErr := f.Close()
		if closeErr != nil {
			return nil, fmt.Errorf("failed to build index: %w | Failed to close file: %v", err, closeErr)
		}
		return nil, fmt.Errorf("failed to build index: %w", err)
	}
	return db, nil
}

func (db *DB) buildindex() error {
	if _, err := db.f.Seek(0, 0); err != nil {
		return err
	}

	r := bufio.NewScanner(db.f)

	var offset int64
	for r.Scan() {
		var e entry
		if err := json.Unmarshal(r.Bytes(), &e); err != nil {
			return err
		}

		lineLen := int64(len(r.Bytes()) + 1)

		if e.V == nil {
			delete(db.index, e.K)
		} else {
			db.index[e.K] = offset
		}
		offset += lineLen
	}
	return r.Err()
}

func (db *DB) Put(k string, v []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	vEnc := base64.StdEncoding.EncodeToString(v)
	e := entry{K: k, V: &vEnc}
	data, _ := json.Marshal(e)
	data = append(data, '\n')

	if _, err := db.f.Write(data); err != nil {
		return err
	}
	db.index[k] = 0 // will be rebuilt on next open leave it for now
	return nil
}

func (db *DB) Get(k string) ([]byte, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.index[k]; !ok {
		return nil, ErrKeyNotFound
	}

	if _, err := db.f.Seek(0, 0); err != nil {
		return nil, err
	}

	r := bufio.NewScanner(db.f)
	for r.Scan() {
		var e entry
		_ = json.Unmarshal(r.Bytes(), &e)
		if e.K != k {
			continue
		}
		if e.V == nil {
			return nil, ErrKeyNotFound
		}
		return base64.StdEncoding.DecodeString(*e.V)
	}
	return nil, ErrKeyNotFound
}

func (db *DB) Delete(k string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.index[k]; !ok {
		return nil
	}

	e := entry{K: k, V: nil}
	data, _ := json.Marshal(e)
	data = append(data, '\n')

	if _, err := db.f.Write(data); err != nil {
		return err
	}
	delete(db.index, k)
	return nil
}

func (db *DB) Close() error {
	return db.f.Close()
}
