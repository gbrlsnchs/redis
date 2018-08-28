package redis_test

import (
	"fmt"
	"testing"

	. "github.com/gbrlsnchs/redis"
)

func TestTxExec(t *testing.T) {
	_, _ = db.Send("FLUSHDB")
	testCases := []struct {
		cmd  string
		args []interface{}
		resp interface{}
		err  error
	}{
		{"GET", []interface{}{"TestTxExec"}, nil, ErrNoResult},
		{"SET", []interface{}{"TestTxExec", "foobar"}, "OK", nil},
		{"GET", []interface{}{"TestTxExec"}, "foobar", nil},
		{"DEL", []interface{}{"TestTxExec"}, 1, nil},
		{"GET", []interface{}{"TestTxExec"}, nil, ErrNoResult},
		{"INCR", []interface{}{"TestTxExec"}, 1, nil},
		{"INCR", []interface{}{"TestTxExec"}, 2, nil},
		{"DECR", []interface{}{"TestTxExec"}, 1, nil},
		{"ECHO", []interface{}{"TestTxExec"}, "TestTxExec", nil},
	}
	for _, tc := range testCases {
		args := fmt.Sprintf("%v", tc.args)
		t.Run(fmt.Sprintf("%s %s", tc.cmd, args), func(t *testing.T) {
			tx, err := db.Multi()
			if want, got := (error)(nil), err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			r, err := tx.Queue(tc.cmd, tc.args...)
			if want, got := (error)(nil), err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if want, got := "QUEUED", r.String(); want != got {
				t.Errorf("want %s, got %s", want, got)
			}
			if r != nil {
				r, err = tx.Exec()
				if want, got := tc.err, err; want != got {
					t.Errorf("want %v, got %v", want, got)
				}
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
				}
			}
		})
	}
}

func TestTxDiscard(t *testing.T) {
	_, _ = db.Send("FLUSHDB")
	testCases := []struct {
		cmd  string
		args []interface{}
	}{
		{"SET", []interface{}{"TestTxDiscard", 12345}},
		{"DEL", []interface{}{"TestTxDiscard"}},
		{"INCR", []interface{}{"TestTxDiscard"}},
		{"DECR", []interface{}{"TestTxDiscard"}},
	}
	const golden = int64(5)
	r, err := db.Send("SET", "TestTxDiscard", golden)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r.String())
	for _, tc := range testCases {
		args := fmt.Sprintf("%v", tc.args)
		t.Run(fmt.Sprintf("%s %s", tc.cmd, args), func(t *testing.T) {
			tx, err := db.Multi()
			if want, got := (error)(nil), err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			r, err := tx.Queue(tc.cmd, tc.args...)
			if want, got := (error)(nil), err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if want, got := "QUEUED", r.String(); want != got {
				t.Errorf("want %s, got %s", want, got)
			}
			r, err = tx.Discard()
			if want, got := (error)(nil), err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if want, got := "OK", r.String(); want != got {
				t.Errorf("want %s, got %s", want, got)
			}

			r, err = db.Send("GET", "TestTxDiscard")
			if want, got := (error)(nil), err; want != got {
				t.Errorf("want %v, got %v", want, got)
			}
			if r != nil {
				if want, got := golden, r.Int64(); want != got {
					t.Errorf("want %d, got %d", want, got)
				}
			}
		})
	}
}
