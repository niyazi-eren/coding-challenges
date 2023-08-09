package server

import (
	"ccwc/redis_server/resp"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ConnHost = "localhost"
	ConnType = "tcp"
)

const (
	SET    = "SET"
	GET    = "GET"
	EXISTS = "EXISTS"
	DEL    = "DEL"
)

const (
	EX   = "EX"   // seconds
	PX   = "PX"   // ms
	EXAT = "EXAT" // unix timeout seconds
	PXAT = "PXAT" // unix timeout ms
)

type RedisValue struct {
	value string
	exp   Expiration
}

type Expiration struct {
	time    string
	option  string
	timeout string
}

type Server struct {
	dict map[string]RedisValue
	mu   sync.Mutex
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
		dict: make(map[string]RedisValue),
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
		reply, err = s.handleSet(reqArgs)
		if err != nil {
			panic(err)
		}
	case GET:
		reply, err = s.handleGet(reqArgs)
		if err != nil {
			panic(err)
		}
	case EXISTS:
		reply = s.handleExists(reqArgs)
	case DEL:
		reply = s.handleDelete(reqArgs)
	}

	_, err = conn.Write([]byte(reply)) // write back the response
	if err != nil {
		fmt.Println("error writing response: " + reply)
	}
	conn.Close() // Close the connection when done
}

// returns the number of keys deleted as a resp integer
func (s *Server) handleDelete(args []string) string {
	count := 0
	for i := 1; i < len(args); i++ {
		key := args[i]
		if _, exists := s.dict[key]; exists {
			delete(s.dict, key)
			count++
		}
	}
	return fmt.Sprintf("%s%d%s", resp.Integers, count, resp.CRLF)
}

// returns the count of existing keys as a resp integer
func (s *Server) handleExists(args []string) string {
	count := 0
	for i := 1; i < len(args); i++ {
		if _, exists := s.dict[args[i]]; exists {
			count++
		}
	}
	return fmt.Sprintf("%s%d%s", resp.Integers, count, resp.CRLF)
}

// returns simple string for OK
// returns bulk string for old value
func (s *Server) handleSet(args []string) (string, error) {
	key := args[1]
	redisValue := RedisValue{value: args[2]}

	err := setExpiration(args, &redisValue)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	oldValue, ok := s.dict[key]
	s.dict[key] = redisValue
	s.mu.Unlock()
	if ok {
		sb := strings.Builder{}
		resp.WriteBulkString(oldValue.value, &sb)
		return sb.String(), nil // if old value is present we return it
	} else {
		return resp.OK, nil
	}
}

func setExpiration(args []string, redisValue *RedisValue) error {
	if len(args) == 5 {
		exp := Expiration{
			option:  args[3],
			timeout: args[4],
			time:    strconv.FormatInt(time.Now().Unix(), 10),
		}
		if isValidExpiration(exp) {
			redisValue.exp = exp
		} else {
			return errors.New("error, unexpected expiration: {opt:" + exp.option + ", timeout:" + exp.timeout + "}")
		}
	}
	return nil
}

// bulk string reply
func (s *Server) handleGet(args []string) (string, error) {
	key := args[1]
	s.mu.Lock()
	val, ok := s.dict[key]
	defer s.mu.Unlock()
	if !ok {
		return resp.NullBulkString, nil
	}
	// check if expired
	if val.exp.timeout != "" {
		expired, err := isExpired(val)
		if err != nil {
			return "", err
		}
		if expired {
			delete(s.dict, key)
			return resp.NullBulkString, nil
		}
	}

	sb := strings.Builder{}
	resp.WriteBulkString(val.value, &sb)
	return sb.String(), nil

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

func isValidExpiration(exp Expiration) bool {
	// Check if the option is valid
	switch exp.option {
	case EX, PX, EXAT, PXAT:
		// Option is valid, check if the timeout is an integer
		_, err := strconv.Atoi(exp.timeout)
		return err == nil
	default:
		// Invalid option
		return false
	}
}

func isExpired(val RedisValue) (bool, error) {
	timestamp, err := strconv.ParseInt(val.exp.time, 10, 64)
	if err != nil {
		return false, fmt.Errorf("error parsing timestamp: %s", err)
	}

	var expirationTime time.Time
	switch val.exp.option {
	case EX:
		timeout, err := strconv.ParseInt(val.exp.timeout, 10, 64)
		if err != nil {
			return false, fmt.Errorf("error parsing timeout: %s", err)
		}
		expirationTime = time.Unix(timestamp, 0).Add(time.Duration(timeout) * time.Second)
	case PX:
		timeout, err := strconv.ParseInt(val.exp.timeout, 10, 64)
		if err != nil {
			return false, fmt.Errorf("error parsing timeout: %s", err)
		}
		expirationTime = time.Unix(timestamp, 0).Add(time.Duration(timeout) * time.Millisecond)
	case EXAT:
		expirationTime = time.Unix(timestamp, 0)
	case PXAT:
		expirationTime = time.Unix(0, timestamp*int64(time.Millisecond))
	default:
		return false, nil
	}

	now := time.Now().Local()
	return now.After(expirationTime), nil
}
