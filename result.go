package redis

import (
	"net"
	"time"

	"github.com/gbrlsnchs/redis/internal"
)

var ErrNoResult = internal.ErrNoResult

type Result struct {
	i      int      // current item being accessed
	null   bool     // null according to RESP
	values [][]byte // raw RESP value
	length int      // cache values len
}

func read(conn net.Conn, times int) (*Result, error) {
	rr := internal.NewReader(conn)
	values, err := rr.ReadN(times)
	if err != nil {
		return nil, err
	}
	return &Result{values: values, length: len(values)}, nil
}

func (r *Result) Bool() bool {
	if r.length == 0 {
		return false
	}
	return Value(r.values[r.index()]).Bool()
}

func (r *Result) Float64() float64 {
	if r.length == 0 {
		return 0
	}
	return Value(r.values[r.index()]).Float64()
}

func (r *Result) Int64() int64 {
	if r.length == 0 {
		return 0
	}
	return Value(r.values[r.index()]).Int64()
}

func (r *Result) IsOK() bool {
	return r.String() == "OK"
}

func (r *Result) Range(fn func(Value)) {
	for _, v := range r.values {
		fn(Value(v))
	}
}

func (r *Result) String() string {
	if r.length == 0 {
		return ""
	}
	return Value(r.values[r.index()]).String()
}

func (r *Result) Unix() time.Time {
	if r.length == 0 {
		return time.Time{}
	}
	i := r.index()
	var nsec int64
	if r.i > 0 { // if there's a next value, it might be microseconds of the unix time, according to RESP
		nsec = Value(r.values[r.index()]).Int64()
	}
	return Value(r.values[i]).Unix(nsec * int64(time.Microsecond))
}

func (r *Result) index() int {
	if r.i >= r.length {
		r.i = 0
	}
	i := r.i
	r.i++
	return i
}
