package storage_test

import (
	"testing"

	"github.com/CinematicCow/lumora/internal/storage"
)

func TestIndexManager_CRUD(t *testing.T) {
	tempDir := t.TempDir()

	im, err := storage.NewIndexManager(tempDir)
	if err != nil {
		t.Fatalf("NewIndexManager failed: %v", err)
	}
	defer im.Close()

	key := "testkey"
	entry := storage.IndexEntry{
		Offset:    1234,
		KeySize:   7,
		ValueSize: 9,
		Timestamp: 123456789,
	}

	// test put
	err = im.PutEntry(key, entry)
	if err != nil {
		t.Fatalf("PutEntry failed: %v", err)
	}

	// test get
	retrieved, err := im.GetEntry(key)
	if err != nil {
		t.Fatalf("GetEntry failed: %v", err)
	}

	if retrieved != entry {
		t.Errorf("Entry mismatch: got %+v, want %+v", retrieved, entry)
	}
}

func TestIndexManager_Persistence(t *testing.T) {
	tempDir := t.TempDir()
	key := "persisted"

	// create and populate
	im1, err := storage.NewIndexManager(tempDir)
	if err != nil {
		t.Fatalf("NewIndexManager1 failed: %v", err)
	}

	entry := storage.IndexEntry{
		Offset: 123,
	}
	im1.PutEntry(key, entry)

	err = im1.Save()
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	im1.Close()

	// reopen and verify
	im2, err := storage.NewIndexManager(tempDir)
	if err != nil {
		t.Fatalf("NewIndexManager2 failed: %v", err)
	}
	defer im2.Close()

	retrieved, err := im2.GetEntry(key)
	if err != nil {
		t.Fatalf("NewIndexManager2_GetEntry failed: %v", err)
	}

	if retrieved != entry {
		t.Errorf("Persistence failed: got %+v, want %+v", retrieved, entry)
	}

}
