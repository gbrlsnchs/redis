package internal

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

func Parse(cmd []byte, args ...interface{}) []byte {
	parsed := make([][]byte, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case string:
			parsed = append(parsed, []byte(v))
		case []byte:
			parsed = append(parsed, v)
		case byte:
			parsed = append(parsed, []byte{v})
		case int:
			parsed = append(parsed, []byte(strconv.Itoa(v)))
		case time.Time:
			t := int(v.Unix())
			parsed = append(parsed, []byte(strconv.Itoa(t)))
		}
	}
	var buf bytes.Buffer
	// Parse the command according to the Redis protocol.
	// Example:
	// 	*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n
	buf.WriteString(fmt.Sprintf("*%d\r\n$%d\r\n%s\r\n", len(parsed)+1, len(cmd), cmd))
	for _, arg := range parsed {
		buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	return buf.Bytes()
}
