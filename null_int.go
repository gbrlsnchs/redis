package redis

// NullInt64 represents an int64 that may be null.
type NullInt64 struct {
	Int64 int64
	Valid bool
}
