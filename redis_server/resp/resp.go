package resp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	CRLF         = "\r\n"
	SimpleString = "+"
	Errors       = "-"
	Integers     = ":"
	BulkStrings  = "$"
	Arrays       = "*"

	NullBulkString = "$-1\r\n"
	OK             = "+OK\r\n"
)

var StringErr = errors.New("string cannot contain a LF or CR")
var TermErr = errors.New("unexpected termination")
var TokenErr = errors.New("unexpected token")
var BytesLenDecodeErr = errors.New("error decoding bulk string length")
var BytesLenExceededErr = errors.New("error the string size cannot be larger than 512MB")
var IncrErr = errors.New("error the value is not an integer or out of range")
var NotAListErr = errors.New("error the value is not a list")

// Encode the command with the RESP protocol
// a command is a RESP Array consisting of only Bulk Strings
func Encode(value string) (string, error) {
	sb := strings.Builder{}

	// Write the Arrays identifier
	sb.WriteString(Arrays)

	tokens := strings.Split(value, " ")
	size := len(tokens)

	// Write the number of tokens as a string followed by CRLF
	sb.WriteString(strconv.Itoa(size))
	sb.WriteString(CRLF)

	// Write each token in bulk strings format
	for i := 0; i < size; i++ {
		WriteBulkString(tokens[i], &sb)
	}
	return sb.String(), nil
}

func WriteBulkString(token string, sb *strings.Builder) {
	sizeToken := len(token)
	sb.WriteString(BulkStrings)
	sb.WriteString(strconv.Itoa(sizeToken))
	sb.WriteString(CRLF)
	sb.WriteString(token)
	sb.WriteString(CRLF)
}

func WriteRespError(msg string) string {
	return fmt.Sprintf("%s%s%s", Errors, msg, CRLF)
}

func WriteRespInt(value int) string {
	return fmt.Sprintf("%s%d%s", Integers, value, CRLF)
}

func Decode(value []byte) (any, error) {
	if len(value) <= 2 {
		return nil, TokenErr
	}
	val := byteSliceToString(value)
	if !strings.HasSuffix(val, CRLF) {
		return nil, TermErr
	}

	dataType := string(value[0])

	data := value[1 : len(val)-2]
	switch dataType {
	case SimpleString:
		ss := string(data)
		if strings.ContainsAny(ss, "\r\n") {
			return nil, StringErr
		}
		return ss, nil
	case Integers:
		return strconv.Atoi(string(data))
	case Errors:
		return nil, errors.New(string(data))
	case BulkStrings:
		return decodeBulkString(data)
	case Arrays:
		return decodeArray(data)
	default:
		return nil, errors.New("unknown data type symbol: " + dataType)
	}
}

func decodeBulkString(data []byte) (any, error) {
	// null bulk string case
	if len(data) <= 2 {
		bytesLen, _ := strconv.Atoi(string(data))
		if bytesLen == -1 {
			return nil, nil
		} else {
			return "", TokenErr
		}
	}

	r := bytes.NewReader(data)
	bytesLen, err := decodeNumber(r)
	if err != nil {
		return "", err
	}

	// cannot be larger than 512MB
	if bytesLen >= 5.12e+8 {
		return "", BytesLenExceededErr
	}

	bsBuf := bytes.Buffer{}
	// write string to buffer
	for i := 0; i < bytesLen; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}
		bsBuf.WriteByte(b)
	}
	return bsBuf.String(), nil
}

func decodeNumber(r *bytes.Reader) (int, error) {
	expectLF := false
	lenBuf := bytes.Buffer{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		char := string(b)
		if expectLF {
			if char != "\n" {
				return 0, BytesLenDecodeErr
			} else {
				break
			}
		} else if char == "\r" {
			expectLF = true
		} else {
			lenBuf.WriteByte(b)
		}
	}
	return strconv.Atoi(lenBuf.String())
}

func decodeArray(data []byte) ([]any, error) {
	arr := make([]any, 0)
	// case empty array
	if len(data) == 1 {
		if string(data) == "0" {
			return arr, nil
		} else {
			return nil, TokenErr
		}
	}

	r := bytes.NewReader(data)
	size, err := decodeNumber(r)
	if err != nil {
		return nil, err
	}

	elemsBuf := bytes.Buffer{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		elemsBuf.WriteString(string(b))
	}

	elems, err := merge(strings.Split(elemsBuf.String(), CRLF))
	if err != nil {
		return nil, err
	}

	for i := 0; i < size; i++ {
		value, err := Decode([]byte(elems[i] + CRLF))
		if err != nil {
			return nil, err
		}
		arr = append(arr, value)
	}
	return arr, nil
}

// merge elements that belong together
func merge(elems []string) ([]string, error) {
	for i := 0; i < len(elems)-1; i++ {
		respType := string(elems[i][0])
		switch respType {
		case BulkStrings:
			if err := mergeBulkString(elems, i); err != nil {
				return nil, err
			}
		case Arrays:
			if err := mergeArray(elems, i); err != nil {
				return nil, err
			}
		}
	}
	return elems, nil
}

func mergeBulkString(elems []string, i int) error {
	var bulk = elems[i] + CRLF + elems[i+1]
	elems[i] = bulk
	elems = append(elems[:i+1], elems[i+2:]...)
	return nil
}

func mergeArray(elems []string, i int) error {
	// read number of elements
	size, err := strconv.Atoi(elems[i][1:])
	if err != nil {
		return err
	}

	values := elems[i] + CRLF
	for j := i + 1; j <= size+i; j++ {
		values += elems[j] + CRLF
	}
	elems[i] = values
	elems = append(elems[:i+1], elems[i+size+1:]...)
	return nil
}

func byteSliceToString(s []byte) string {
	n := bytes.IndexByte(s, 0)
	if n >= 0 {
		s = s[:n]
	}
	return string(s)
}
