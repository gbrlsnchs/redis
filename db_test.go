package redis_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	. "github.com/gbrlsnchs/redis"
)

var db *DB

func TestMain(m *testing.M) {
	var err error
	db, err = Open("localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	_, _ = db.Send("FLUSHDB")
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	testCases := []struct {
		args     []interface{}
		expected string
	}{
		{nil, "PONG"},
		{[]interface{}{"hello"}, "hello"},
		{[]interface{}{"hello world"}, "hello world"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.args), func(t *testing.T) {
			msg, err := db.Ping(tc.args...)
			if want, got := (error)(nil), err; want != got {
				t.Errorf("want %t, got %t", want, got)
			}
			if want, got := tc.expected, msg; want != got {
				t.Errorf("want %s, got %s", want, got)
			}
		})
	}
}

func TestSend(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		cmd  string
		args []interface{}
		resp interface{}
		err  error
	}{
		{"GET", []interface{}{"TestSend"}, nil, ErrNoResult},
		{"SET", []interface{}{"TestSend", "foobar"}, "OK", nil},
		{"GET", []interface{}{"TestSend"}, "foobar", nil},
		{"DEL", []interface{}{"TestSend"}, 1, nil},
		{"GET", []interface{}{"TestSend"}, nil, ErrNoResult},
		{"INCR", []interface{}{"TestSend"}, 1, nil},
		{"INCR", []interface{}{"TestSend"}, 2, nil},
		{"DECR", []interface{}{"TestSend"}, 1, nil},
		{"ECHO", []interface{}{"TestSend"}, "TestSend", nil},
		{"TIME", nil, now, nil},
	}
	for _, tc := range testCases {
		args := fmt.Sprintf("%v", tc.args)
		t.Run(fmt.Sprintf("%s %s", tc.cmd, args), func(t *testing.T) {
			r, err := db.Send(tc.cmd, tc.args...)
			if want, got := tc.err, err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if r != nil {
				switch vv := tc.resp.(type) {
				case bool:
					if want, got := vv, r.Bool(); want != got {
						t.Errorf("want %t, got %t", want, got)
					}
				case int:
					if want, got := vv, int(r.Int64()); want != got {
						t.Errorf("want %d, got %d", want, got)
					}
				case string:
					s := r.String()
					if want, got := len(vv), len(s); want != got {
						t.Errorf("want %d, got %d", want, got)
					}
					if want, got := vv, s; want != got {
						t.Errorf("want %s, got %s", want, got)
					}
				case time.Time:
					rtime := r.Unix()
					if want, got := true, vv.Before(rtime); want != got {
						t.Errorf("want %s to be before %s", vv.String(), rtime.String())
					}
				}
			}
		})
	}
}
