package redis

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"net"
	"strings"
	"sync"
)

const (
	Port = ":16379"
)

/*
Attributes of the server
	- host
	- port (default: 16379)
	- Key-value store (a dict)
	- A protocol object
	- run_server()
	- conn_handler()
	- get_resp()
*/

type MemoryManager struct {
	MemMap map[string]DynamicDataStruct
	mtx    sync.RWMutex
}

func (m *MemoryManager) GetHanlder(key string) (response ResponseProtocol) {
	m.mtx.RLock()
	data, ok := m.MemMap[key]
	m.mtx.RUnlock()

	if ok {
		response.Data = DataStruct{key, data}
		response.Error = ""
	} else {
		response.Data = EmptyDataStruct
		response.Error = "Not Found"
	}
	return response
}

func (m *MemoryManager) SetHandler(data DataStruct) (response ResponseProtocol) {
	m.mtx.Lock()
	m.MemMap[data.Key] = data.Data
	m.mtx.Unlock()

	response.Error = ""
	response.Data = EmptyDataStruct

	//log.Println(m.MemMap)
	return response
}

func (m *MemoryManager) DeleteHandler(key string) (response ResponseProtocol) {
	m.mtx.Lock()
	if _, ok := m.MemMap[key]; ok {
		delete(m.MemMap, key)
	}
	m.mtx.Unlock()

	response.Error = ""
	response.Data = EmptyDataStruct
	return response
}

func (m *MemoryManager) EvictHandler() (response ResponseProtocol) {
	m.mtx.Lock()
	m.MemMap = make(map[string]DynamicDataStruct)
	m.mtx.Unlock()

	response.Error = ""
	response.Data = EmptyDataStruct
	return response
}

var MMObject MemoryManager

// HandleFunc is a function that handles an incoming command.
// It receives the open connection wrapped in a `ReadWriter` interface.
type HandleFunc func(*bufio.ReadWriter)

/* Endpoint provides an endpoint to other processess
that they can send data to. */
type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc

	// Maps are not thread-safe, so we need a mutex to control access.
	m sync.RWMutex
}

/* NewEndpoint creates a new endpoint. Too keep things simple,
the endpoint listens on a fixed port number. */
func NewEndpoint() *Endpoint {
	// Create a new Endpoint with an empty list of handler funcs.
	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

/* AddHandleFunc adds a new function for handling incoming data. */
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

/* Listen starts listening on the endpoint port on all interfaces.
At least one handler function must have been added
through AddHandleFunc() before. */
func (e *Endpoint) Listen() error {
	var err error
	e.listener, err = net.Listen("tcp", Port)
	if err != nil {
		return errors.Wrapf(err, "Unable to listen on port %s\n", Port)
	}
	log.Println("Listen on", e.listener.Addr().String())
	for {
		log.Println("Accept a connection request.")
		conn, err := e.listener.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go e.handleMessages(conn)
	}
}

/* handleMessages reads the connection up to the first newline.
Based on this string, it calls the appropriate HandleFunc. */
func (e *Endpoint) handleMessages(conn net.Conn) {
	defer conn.Close()

	decoder := json.NewDecoder(conn)
	var request RequestProtocol
	if err := decoder.Decode(&request); err != nil {
		log.Println("Error while decoding: ", err)
		return
	}
	var response ResponseProtocol
	log.Println(request.Command)
	switch request.Command {
	case "GET":
		response = MMObject.GetHanlder(request.Data.Key)
	case "SET":
		response = MMObject.SetHandler(request.Data)
	case "DEL":
		response = MMObject.DeleteHandler(request.Data.Key)
	case "EVICT":
		response = MMObject.EvictHandler()
	default:
		response.Error = "Unknown Command"
		response.Data = EmptyDataStruct
	}
	encoder := json.NewEncoder(conn)
	encoder.Encode(response)
}

func handleRequests(rw *bufio.ReadWriter) {
	log.Print("Receive STRING message:")
	s, err := rw.ReadString('\n')
	if err != nil {
		log.Println("Cannot read from connection.\n", err)
	}
	s = strings.Trim(s, "\n ")
	log.Println(s)
	_, err = rw.WriteString("Thank you.\n")
	if err != nil {
		log.Println("Cannot write to connection.\n", err)
	}
	err = rw.Flush()
	if err != nil {
		log.Println("Flush failed.", err)
	}
}

/* server listens for incoming requests and dispatches them to
registered handler functions. */
func ServerFunc() error {
	MMObject.MemMap = make(map[string]DynamicDataStruct)
	endpoint := NewEndpoint()

	// Add the handle funcs.
	endpoint.AddHandleFunc("*", handleRequests)

	// Start listening.
	return endpoint.Listen()
}
