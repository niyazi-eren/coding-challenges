package main

import (
	"ccwc/redis_server/resp"
	"fmt"
	"net"
	"os"
	"reflect"
	"strings"
	"sync"
)

const (
	ConnHost = "localhost"
	ConnType = "tcp"
)

const (
	SET = "SET"
	GET = "GET"
)

type Server struct {
	dict map[string]string
	mu   sync.Mutex
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
		dict: make(map[string]string),
	}
}

func (s *Server) Run() {
	// Listen for incoming connections.
	l, err := net.Listen(ConnType, ConnHost+":"+s.port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	fmt.Println("Listening on " + ConnHost + ":" + s.port)

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections concurrently
		go s.handleRequest(conn)
	}
}

func (s *Server) Close() {
	s.Close()
}

func (s *Server) handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	// read the incoming connection into the buffer
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}

	// decode request
	req, err := resp.Decode(buf)
	if err != nil {
		fmt.Println("Error decoding requests:", err.Error())
	}

	reqArgs, err := anyToStringArray(req) // convert the request to an array of string
	if err != nil {
		fmt.Println("error converting request args")
	}

	cmd := reqArgs[0]
	var reply string
	switch cmd {
	case SET:
		reply = s.handleSet(reqArgs)
	case GET:
		reply = s.handleGet(reqArgs)
	}

	_, err = conn.Write([]byte(reply)) // write back the response
	if err != nil {
		fmt.Println("error writing response: " + reply)
	}
	conn.Close() // Close the connection when  done
}

// returns simple string for OK
// returns bulk string for old value
func (s *Server) handleSet(reqArgs []string) string {
	key := reqArgs[1]
	value := reqArgs[2]
	s.mu.Lock()
	oldValue, ok := s.dict[key]
	s.dict[key] = value
	s.mu.Unlock()
	if ok {
		sb := strings.Builder{}
		resp.WriteBulkString(oldValue, &sb)
		return sb.String() // if old value is present we return it
	} else {
		return resp.OK
	}
}

// bulk string reply
func (s *Server) handleGet(reqArgs []string) string {
	key := reqArgs[1]
	s.mu.Lock()
	val, ok := s.dict[key]
	defer s.mu.Unlock()
	if !ok {
		return resp.NullBulkString
	} else {
		sb := strings.Builder{}
		resp.WriteBulkString(val, &sb)
		return sb.String()
	}
}

// helper method to convert an any value to a string array
func anyToStringArray(value any) ([]string, error) {
	// Check if the value is an array or slice
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input value is not an array or slice")
	}
	// Convert each element of the array to a string and create a new array
	var result []string
	for i := 0; i < val.Len(); i++ {
		elem := fmt.Sprint(val.Index(i).Interface())
		result = append(result, elem)
	}
	return result, nil
}
