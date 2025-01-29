package storage_test

import (
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

			if string(readRecord.Value) != tt.key {
				t.Errorf("Value mismatch: got %q, want %q", readRecord.Value, tt.value)
			}

		})
	}
}
