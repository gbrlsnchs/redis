package redis

type Type byte

const (
	StringType     Type = '+'
	ErrorType      Type = '-'
	IntType        Type = ':'
	BulkStringType Type = '$'
	ArrayType      Type = '*'
)

func (t Type) String() string {
	switch t {
	case StringType:
		return "(string)"

	case ErrorType:
		return "(error)"

	case IntType:
		return "(integer)"

	case BulkStringType:
		return "(bulk string)"

	case ArrayType:
		return "(array)"

	default:
		return "(invalid type)"
	}
}
