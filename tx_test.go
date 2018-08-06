package redis_test

import (
	"testing"

	. "github.com/gbrlsnchs/redis"
)

func TestTxCommit(t *testing.T) {
	tx, err := db.Begin()
	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
		t.FailNow()
	}

	err = tx.Exec("SET", "commit", 1)
	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
		t.FailNow()
	}

	err = tx.Commit()
	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
	}

	var commit NullInt64
	err = db.QueryRow("GET", "commit").Scan(&commit)
	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
	}
	if want, got := true, commit.Valid; want != got {
		t.Errorf("want %t, got %t\n", want, got)
		t.FailNow()
	}
	if want, got := int64(1), commit.Int64; want != got {
		t.Errorf("want %d, got %d\n", want, got)
	}
}

func TestTxRollback(t *testing.T) {
	tx, err := db.Begin()
	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
		t.FailNow()
	}
	err = tx.Exec("SET", "rb", 1)
	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
		t.FailNow()
	}
	var rb NullInt64
	err = tx.Rollback()
	if want, got := (error)(nil), err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
		t.FailNow()
	}
	err = db.QueryRow("GET", "rb").Scan(&rb)
	if want, got := ErrNoResult, err; want != got {
		t.Errorf("want %v, got %v\n", want, got)
	}
	if want, got := false, rb.Valid; want != got {
		t.Errorf("want %t, got %t\n", want, got)
		t.FailNow()
	}
}
