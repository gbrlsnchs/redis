package redis_test

import (
	"log"
	"os"
	"testing"

	. "github.com/gbrlsnchs/redis"
	. "github.com/gbrlsnchs/redis/internal"
)

var db *DB

func TestMain(m *testing.M) {
	var err error
	db, err = Open("localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestErrNoResult(t *testing.T) {
	var dump Dump
	err := db.QueryRow("GET", "never_set").Scan(&dump)
	if want, got := ErrNoResult, err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
	}
}

func TestPing(t *testing.T) {
	pong, err := db.Ping()
	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
		t.FailNow()
	}
	if want, got := "PONG", pong; want != got {
		t.Errorf("want %s, got %s\n", want, got)
	}

	msg := "Hello Word"
	pong, err = db.Ping(msg)
	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %s, got %s\n", want, got)
		t.FailNow()
	}
	if want, got := msg, pong; want != got {
		t.Errorf("want %s, got %s\n", want, got)
	}
}
