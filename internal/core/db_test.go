package core_test

import (
	"fmt"
	"testing"

	"github.com/CinematicCow/lumora/internal/core"
	"github.com/CinematicCow/lumora/internal/storage"
)

func TestDB_BasicOps(t *testing.T) {
	tempDir := t.TempDir()

	db, err := core.Open(tempDir)

	if err != nil {
		t.Fatalf("Open failed: %v", err)

	}

	key := "testKey"
	value := []byte("testValue")

	// test put
	err = db.Put(key, value)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	// test get
	retrieved, err := db.Get(key)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	if string(retrieved) != string(value) {
		t.Errorf("Value mismatch: got %q, want %q", string(retrieved), string(value))
	}

	// test non existent key
	_, err = db.Get("brrr")
	if err != storage.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}

func TestDB_Concurrency(t *testing.T) {
	tempDir := t.TempDir()

	db, _ := core.Open(tempDir)

	key := "counter"
	iterations := 100

	// concurrent writes
	t.Run("concurrent writes", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			go func(n int) {
				db.Put(key, []byte(fmt.Sprintf("%d", n)))
			}(i)
		}
	})

	// concurrent reads
	t.Run("concurrent reads", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			go func() {

				_, err := db.Get(key)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}()
		}
	})
}

func TestDB_Reopen(t *testing.T) {
	tempDir := t.TempDir()

	// first session
	db1, _ := core.Open(tempDir)
	db1.Put("persistent", []byte("zaza"))
	db1.Close()

	// reopen
	db2, _ := core.Open(tempDir)
	defer db2.Close()

	val, err := db2.Get("persistent")
	if err != nil {
		t.Fatalf("db2 get fail: %q", err)
	}

	if string(val) != "zaza" {
		t.Errorf("Persistence failed: got %q, want %q", val, "zaza")
	}

}

func TestDB_Delete(t *testing.T) {
	tempDir := t.TempDir()

	db, err := core.Open(tempDir)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	key := "yer"
	value := []byte("mom")

	err = db.Put(key, value)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	err = db.Delete(key)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = db.Get(key)
	if err != storage.ErrKeyNotFound {
		t.Errorf("Expected ErrKeyNotFound, got %v", err)
	}
}
