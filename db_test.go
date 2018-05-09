package redis_test

import (
	"log"
	"testing"

	. "github.com/gbrlsnchs/redis"
)

var db *DB

func init() {
	var err error
	db, err = Open("localhost", 6379, 10)

	if err != nil {
		log.Fatal(err)
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

	pong, err = db.Ping("Hello World")

	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %s, got %s\n", want, got)
		t.FailNow()
	}

	if want, got := "Hello World", pong; want != got {
		t.Errorf("want %s, got %s\n", want, got)
	}
}
