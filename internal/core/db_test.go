package core_test

import (
	"testing"

	"github.com/CinematicCow/lumora/internal/core"
)

func TestDB_BasicOps(t *testing.T) {
	dir := t.TempDir()
	db, err := core.Open(dir)
	if err != nil {
		t.Fatal(err)
	}

	defer func ()  {
		if err := db.Close(); err!=nil {
			t.Fatalf("close failed: %v",err)
		}
	}()

	key, val := "k", []byte("v")
	if err := db.Put(key, val); err != nil {
		t.Fatal(err)
	}

	res, err := db.Get(key)
	if err != nil {
		t.Fatal(err)
	}
	if string(res) != string(val) {
		t.Fatalf("got %q want %q", res, val)
	}

	if err := db.Delete(key); err != nil {
		t.Fatal(err)
	}
	if _, err := db.Get(key); err != core.ErrKeyNotFound {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestDB_Reopen(t *testing.T) {
	dir := t.TempDir()

	db1, _ := core.Open(dir)
	_ = db1.Put("yo", []byte("mama"))
	_ = db1.Close()

	db2, _ := core.Open(dir)

	defer func ()  {
		if err := db2.Close(); err!=nil {
			t.Fatalf("close failed: %v",err)
		}
	}()

	v, err := db2.Get("yo")
	if err != nil || string(v) != "mama" {
		t.Fatalf("reopen failed: %v %q", err, v)
	}
}
