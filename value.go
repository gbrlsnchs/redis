package redis

import (
	"strconv"
	"time"
)

// Value is a raw Redis value.
type Value []byte

// Bool converts a value to a boolean value.
func (v Value) Bool() bool {
	return v.Int64() != 0
}

// Float64 converts a value to a float64.
// If the conversion fails, it returns 0 instead.
func (v Value) Float64() float64 {
	n, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0
	}
	return n
}

// Int64 converts a value to an int64.
// If the conversion fails, it returns 0 instead.
func (v Value) Int64() int64 {
	n, err := strconv.ParseInt(string(v), 10, 64)
	if err != nil {
		return 0
	}
	return n
}

// String converts a value to a string.
func (v Value) String() string {
	return string(v)
}

func (v Value) Uint64() uint64 {
	n, err := strconv.ParseUint(string(v), 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func (v Value) Unix(nsec int64) time.Time {
	n := v.Int64()
	return time.Unix(n, nsec)
}
