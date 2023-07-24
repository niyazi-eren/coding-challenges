package resp_test

import (
	"ccwc/redis_server/resp"
	"errors"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	cases := []struct {
		Description string
		Want        any
		Command     string
	}{
		{"unexpected token error", resp.TokenErr, "+\r"},
		{"unexpected termination error", resp.TermErr, "+hello world\r"},
		{"string hello world is decoded", "hello world", "+hello world\r\n"},
		{"string containing a LF returns an error", resp.StringErr, "+hello wo\nrld\r\n"},
		{"integer 12354111 is decoded", 12354111, ":12354111\r\n"},
		{"integer 0 is decoded", 0, ":0\r\n"},
		{"error 'error message' is decoded", errors.New("error message"), "-error message\r\n"},
		{"bulk string is decoded", "hel\nlo", "$6\r\nhel\nlo\r\n"},
		{"empty string is decoded", "", "$0\r\n\r\n"},
		{"null bulk string returns nil", nil, resp.NullBulkString},
		{"array of int 1, 2, 3 is decoded", []any{1, 2, 3}, "*3\r\n:1\r\n:2\r\n:3\r\n"},
		{"array of any values is decoded", []any{1, 2, 3, 4, "hello"}, "*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$5\r\nhello\r\n"},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			got, err := resp.Decode([]byte(test.Command))
			if err != nil {
				switch test.Want.(type) {
				// case when expecting an error
				case error:
					wantErr := test.Want.(error)
					if err.Error() != wantErr.Error() {
						t.Errorf("got %q, want %q", err, test.Want)
					}
				default:
					t.Errorf("unexpected error: %s", err.Error())
				}
			} else if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("got %q, want %q", got, test.Want)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	cases := []struct {
		Description string
		Want        any
		Command     string
	}{
		{"command LLEN mylist is encoded", "*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n", "LLEN mylist"},
	}

	for _, test := range cases {
		t.Run(test.Description, func(t *testing.T) {
			got, err := resp.Encode(test.Command)
			if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			} else if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("got %q, want %q", got, test.Want)
			}
		})
	}
}
