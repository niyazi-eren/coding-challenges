package server_test

import (
	"ccwc/redis_server"
	"ccwc/redis_server/resp"
	"fmt"
	"net"
	"testing"
	"time"
)

var testPort = ":8888"

func TestServer_Run(t *testing.T) {
	conn, err := net.Dial("tcp", testPort)
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer conn.Close()
}

func TestServer_Set(t *testing.T) {
	tests := []struct {
		cmd  string
		want string
	}{
		{"SET name JOHN", "OK"},
		{"SET name JANE", "JOHN"},
	}

	for _, tt := range tests {
		got, err := send(tt.cmd)
		if err != nil {
			t.Errorf(err.Error())
		}

		if got != tt.want {
			t.Errorf("for command %q, got %q, want %q", tt.cmd, got, tt.want)
		}
	}
}

func TestServer_Get(t *testing.T) {
	_, err := send("SET name JOHN")
	if err != nil {
		t.Errorf(err.Error())
	}

	testsGet := []struct {
		cmd  string
		want any
	}{
		{"GET name", "JOHN"},
		{"GET doesntexist", nil},
	}

	for _, tt := range testsGet {
		got, err := send(tt.cmd)
		if err != nil {
			t.Errorf(err.Error())
		} else {
			if got != tt.want {
				t.Errorf("got %q, want %q", err, tt.want)
			}
		}
	}
}

func TestServer_SetExpire(t *testing.T) {
	_, err := send("SET name JOHN PX 10")
	_, err = send("SET test TEST PX 3500")
	if err != nil {
		t.Errorf(err.Error())
	}

	testsGet := []struct {
		cmd  string
		want any
	}{
		{"GET name", nil},
		{"GET test", "TEST"},
	}

	time.Sleep(time.Millisecond * 10)

	for _, tt := range testsGet {
		got, err := send(tt.cmd)
		if err != nil {
			t.Errorf(err.Error())
		} else {
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		}
	}
}

func TestServer_Exists_Del(t *testing.T) {
	_, err := send("SET name JOHN")
	if err != nil {
		t.Errorf(err.Error())
	}

	tests := []struct {
		cmd  string
		want any
	}{
		{"EXISTS fail", 0},
		{"EXISTS fail name name", 2},
		{"DEL name name", 1},
		{"DEL fail", 0},
		{"EXISTS name", 0},
	}

	for _, tt := range tests {
		got, err := send(tt.cmd)
		if err != nil {
			t.Errorf(err.Error())
		} else {
			if got != tt.want {
				t.Errorf("got %q, want %q", err, tt.want)
			}
		}
	}
}

// executed before every test
func init() {
	s := server.NewServer("8888")
	go s.Run()
}

func send(cmd string) (any, error) {
	respCmd, err := resp.Encode(cmd)
	if err != nil {
		return "", err
	}
	conn, err := net.Dial("tcp", "localhost"+testPort)
	defer conn.Close()
	if err != nil {
		return "", err
	}
	_, err = conn.Write([]byte(respCmd))
	if err != nil {
		return "", err
	}
	buf := make([]byte, 0124)
	_, err = conn.Read(buf)
	if err != nil {
		return "", err
	}

	response, err := resp.Decode(buf)
	if err != nil {
		fmt.Println("error decoding: ", err.Error())
	}
	return response, nil
}
