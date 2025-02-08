package storage_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/CinematicCow/lumora/internal/storage"
)

func TestDataManager_WriteReadRecord(t *testing.T) {
	tempDir := t.TempDir()

	dm, err := storage.NewDataManager(tempDir)
	if err != nil {
		t.Fatal("Failed to create DataManager: ", err)
	}

	tests := []struct {
		name    string
		key     string
		value   []byte
		wantErr bool
	}{
		{"basic record", "testKey", []byte("testValue"), false},
		{"empty value", "empty", []byte{}, false},
		{"large value", "big", make([]byte, 1024*1024), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record := &storage.DataRecord{
				Timestamp: uint32(123456789),
				Key:       []byte(tt.key),
				Value:     tt.value,
				KeySize:   uint32(len(tt.key)),
				ValueSize: uint32(len(tt.value)),
			}

			offset, err := dm.WriteRecord(record)
			if (err != nil) != tt.wantErr {
				t.Fatalf("WriteRecord() error = %v, wantErr %v", err, tt.wantErr)
			}

			readRecord, err := dm.ReadRecord(offset)
			if err != nil {
				t.Fatal("ReadRecord() failed: ", err)
			}

			if string(readRecord.Key) != tt.key {
				t.Errorf("Key mismatch: got %q, want %q", readRecord.Key, tt.key)
			}

			if !bytes.Equal(readRecord.Value, tt.value) {
				t.Errorf("Value mismatch: got %v, want %v", readRecord.Value, tt.value)
			}

		})
	}
}

func TestDataManager_CorruptedFile(t *testing.T) {
	tempDir := t.TempDir()

	dataPath := filepath.Join(tempDir, storage.DataFileName)
	corruptFile, err := os.Create(dataPath)
	if err != nil {
		t.Fatalf("Failed to create data file: %v", err)
	}

	header := make([]byte, 12)
	binary.BigEndian.AppendUint32(header[0:4], 1)
	binary.BigEndian.AppendUint32(header[4:8], 1<<24+1)
	binary.BigEndian.AppendUint32(header[8:12], 0)

	_, err = corruptFile.Write(header)
	if err != nil {
		t.Fatal(err)
	}
	corruptFile.Close()

	// init data manager with corrupted file
	dm, err := storage.NewDataManager(tempDir)
	if err != nil {
		t.Fatalf("DataManager failed: %v", err)
	}
	defer dm.Close()

	_, err = dm.ReadRecord(0)
	if err == nil {
		t.Fatal("Expected error reading corrupted record")
	}

	// verify err type
	if !errors.Is(err, storage.ErrDataCorruption) {
		t.Errorf("Expected ErrDataCorruption, got: %v", err)
	}
}
