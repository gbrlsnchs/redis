package internal

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"time"
	"unsafe"
)

type Writer struct {
	w   io.Writer
	buf *bytes.Buffer
}

func NewWriter(w io.Writer) *Writer {
	var buf bytes.Buffer
	return &Writer{w, &buf}
}

func (w *Writer) WriteCmd(cmd string, args ...interface{}) (int, error) {
	return w.w.Write(w.parse(cmd, args...))
}

func (w *Writer) parse(cmd string, args ...interface{}) []byte {
	// Parse the command according to the Redis protocol.
	// Example:
	// 	*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n
	parsed := w.parseArgs(args...)
	w.buf.WriteString(fmt.Sprintf("*%d\r\n$%d\r\n%s\r\n", len(parsed)+1, len(cmd), cmd))
	for _, arg := range parsed {
		w.buf.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	defer w.buf.Reset()
	return w.buf.Bytes()
}

func (w *Writer) parseArgs(args ...interface{}) []string {
	parsed := make([]string, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		case bool:
			n := 0
			if v {
				n = 1
			}
			parsed = append(parsed, strconv.Itoa(n))
		case string:
			parsed = append(parsed, v)
		case []byte:
			parsed = append(parsed, *(*string)(unsafe.Pointer(&v)))
		case byte:
			parsed = append(parsed, *(*string)(unsafe.Pointer(&v)))
		case int:
			parsed = append(parsed, strconv.FormatInt(int64(v), 10))
		case int8:
			parsed = append(parsed, strconv.FormatInt(int64(v), 10))
		case int16:
			parsed = append(parsed, strconv.FormatInt(int64(v), 10))
		case int32:
			parsed = append(parsed, strconv.FormatInt(int64(v), 10))
		case int64:
			parsed = append(parsed, strconv.FormatInt(v, 10))
		case uint:
			parsed = append(parsed, strconv.FormatUint(uint64(v), 10))
		case uint16:
			parsed = append(parsed, strconv.FormatUint(uint64(v), 10))
		case uint32:
			parsed = append(parsed, strconv.FormatUint(uint64(v), 10))
		case uint64:
			parsed = append(parsed, strconv.FormatUint(v, 10))
		case float32:
			parsed = append(parsed, strconv.FormatFloat(float64(v), 'f', -1, 64))
		case float64:
			parsed = append(parsed, strconv.FormatFloat(v, 'f', -1, 64))
		case time.Time:
			t := int(v.Unix())
			parsed = append(parsed, strconv.Itoa(t))
		}
	}
	return parsed
}
