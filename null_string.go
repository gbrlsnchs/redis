package redis

// NullString represents a string that may be null.
type NullString struct {
	String string
	Valid  bool
}
