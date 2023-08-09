package server

import (
	"bufio"
	"ccwc/redis_server/resp"
	"encoding/binary"
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
	INCR   = "INCR"
	DECR   = "DECR"
	LPUSH  = "LPUSH"
	RPUSH  = "RPUSH"
	SAVE   = "SAVE"
	LOAD   = "LOAD"
)

const (
	EX   = "EX"   // seconds
	PX   = "PX"   // ms
	EXAT = "EXAT" // unix timeout seconds
	PXAT = "PXAT" // unix timeout ms
)

type RedisValue struct {
	value any
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
	case INCR:
		reply = s.handleIncrDecr(reqArgs, true)
	case DECR:
		reply = s.handleIncrDecr(reqArgs, false)
	case RPUSH:
		reply = s.handleRPush(reqArgs)
	case LPUSH:
		reply = s.handleLPush(reqArgs)
	case SAVE:
		reply = s.handleSave()
	case LOAD:
		reply = s.handleLoad()
	}
	_, err = conn.Write([]byte(reply)) // write back the response
	if err != nil {
		fmt.Println("error writing response: " + reply)
	}
	conn.Close() // Close the connection when done
}

// The SAVE commands performs a synchronous save of the dataset
// producing a point in time snapshot of all the data inside the Redis instance,
// in the form of an RDB file.
func (s *Server) handleSave() string {
	file, err := os.Create("snapshot.rdb")
	if err != nil {
		return resp.WriteRespError(err.Error())
	}
	defer file.Close()

	// encode and write the data
	for k, v := range s.dict {
		value := v.value
		data := k + ":" + anyToString(value) + "\n"
		err = binary.Write(file, binary.LittleEndian, []byte(data))
		if err != nil {
			return resp.WriteRespError("error binary encoding")
		}
	}
	return resp.OK
}

func (s *Server) handleLoad() string {
	file, err := os.Open("snapshot.rdb")
	if err != nil {
		resp.WriteRespError(err.Error())
	}
	defer file.Close()

	m := make(map[string]RedisValue)
	reader := bufio.NewReader(file)
	for {
		// Read a chunk of data from the file
		bytes, err := reader.ReadBytes('\n')
		data := string(bytes)
		if data == "" {
			break
		}
		pair := strings.Split(data, ":")
		key := pair[0]
		value := pair[1]
		m[key] = RedisValue{value: value}
		if err != nil {
			break // Break loop on end of file or error
		}
	}
	s.dict = m
	return resp.OK
}

// Returns Integer reply: the length of the list after the push operations.
func (s *Server) handleLPush(args []string) string {
	key := args[1]
	redisVal, exists := s.dict[key]
	//If key does not exist, it is created as empty list
	if !exists {
		redisVal.value = make([]string, 0)
	}

	arr, err := anyToStringArray(redisVal.value)
	//When key holds a value that is not a list, an error is returned.
	if err != nil {
		resp.WriteRespError(resp.NotAListErr.Error())
	}

	// Insert all the specified values at the head of the list stored at key.
	for i := 2; i < len(args); i++ {
		head := []string{args[i]}
		arr = append(head, arr...)
	}

	redisVal.value = arr
	s.dict[key] = redisVal
	return resp.WriteRespInt(len(arr))
}

// Insert all the specified values at the end of the list stored at key.
// Returns Integer reply: the length of the list after the push operations.
func (s *Server) handleRPush(args []string) string {
	key := args[1]
	redisVal, exists := s.dict[key]
	//If key does not exist, it is created as empty list
	if !exists {
		redisVal.value = make([]string, 0)
	}

	arr, err := anyToStringArray(redisVal.value)
	//When key holds a value that is not a list, an error is returned.
	if err != nil {
		resp.WriteRespError(resp.NotAListErr.Error())
	}

	for i := 2; i < len(args); i++ {
		arr = append(arr, args[i])
	}

	redisVal.value = arr
	s.dict[key] = redisVal
	return resp.WriteRespInt(len(arr))
}

// Return Integer reply: the value of key after the increment or decrement
func (s *Server) handleIncrDecr(args []string, increment bool) string {
	key := args[1]
	_, exists := s.dict[key]

	// If the key does not exist, it is set to 0 before performing the operation
	if !exists {
		redisValue := RedisValue{value: strconv.Itoa(0)}
		s.dict[key] = redisValue
	}

	redisVal, _ := s.dict[key]
	// An error is returned if the key contains a value of the wrong type or contains a string that can not be represented as integer
	val, err := strconv.ParseInt(anyToString(redisVal.value), 10, 64)
	if err != nil {
		return resp.WriteRespError(resp.IncrErr.Error())
	}
	// Increment or decrements the number stored at key
	if increment {
		val++
	} else {
		val--
	}

	redisVal.value = strconv.FormatInt(val, 10)
	s.dict[key] = redisVal
	return fmt.Sprintf("%s%d%s", resp.Integers, redisVal.value, resp.CRLF)
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
	return resp.WriteRespInt(count)
}

// returns the count of existing keys as a resp integer
func (s *Server) handleExists(args []string) string {
	count := 0
	for i := 1; i < len(args); i++ {
		if _, exists := s.dict[args[i]]; exists {
			count++
		}
	}
	return resp.WriteRespInt(count)
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
		resp.WriteBulkString(anyToString(oldValue.value), &sb)
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
	resp.WriteBulkString(anyToString(val.value), &sb)
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

// helper method to convert an any value to a array
func anyToString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		return ""
	}
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
