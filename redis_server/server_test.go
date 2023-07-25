package main

import (
	"ccwc/redis_server/resp"
	"fmt"
	"net"
	"testing"
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
	cmd := "SET name JOHN"
	got, err := send(cmd)
	if err != nil {
		t.Errorf(err.Error())
	}
	want := "OK"
	if got != want {
		t.Errorf("got %q, want %q", err, want)
	}

	cmd = "SET name JANE"
	got, err = send(cmd)
	want = "JOHN"
	if got != want {
		t.Errorf("got %q, want %q", err, want)
	}
}

func TestServer_Get(t *testing.T) {
	_, err := send("SET name JOHN")
	if err != nil {
		t.Errorf(err.Error())
	}

	got, err := send("GET name")
	if err != nil {
		t.Errorf(err.Error())
	}

	want := "JOHN"
	if got != want {
		t.Errorf("got %q, want %q", err, want)
	}

	got, err = send("GET doesntexist")
	if err != nil {
		t.Errorf(err.Error())
	}

	if got != nil {
		t.Errorf("got %q, want %q", err, want)
	}
}

// executed before every test
func init() {
	s := NewServer("8888")
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
