package redis

import (
	"bufio"
	"encoding/gob"
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
type Server struct {



}



type complexData struct {
	N int
	S string
	M map[string]int
	P []byte
	C *complexData
}


// HandleFunc is a function that handles an incoming command.
// It receives the open connection wrapped in a `ReadWriter` interface.
type HandleFunc func(*bufio.ReadWriter)

// Endpoint provides an endpoint to other processess
// that they can send data to.
type Endpoint struct {
	listener net.Listener
	handler  map[string]HandleFunc

	// Maps are not thread-safe, so we need a mutex to control access.
	m sync.RWMutex
}

// NewEndpoint creates a new endpoint. Too keep things simple,
// the endpoint listens on a fixed port number.
func NewEndpoint() *Endpoint {
	// Create a new Endpoint with an empty list of handler funcs.
	return &Endpoint{
		handler: map[string]HandleFunc{},
	}
}

// AddHandleFunc adds a new function for handling incoming data.
func (e *Endpoint) AddHandleFunc(name string, f HandleFunc) {
	e.m.Lock()
	e.handler[name] = f
	e.m.Unlock()
}

// Listen starts listening on the endpoint port on all interfaces.
// At least one handler function must have been added
// through AddHandleFunc() before.
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


type temp_interface struct {
	str string
	num int
}

// handleMessages reads the connection up to the first newline.
// Based on this string, it calls the appropriate HandleFunc.
func (e *Endpoint) handleMessages(conn net.Conn) {
	// Wrap the connection into a buffered reader for easier reading.
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	var t temp_interface
	decoder := json.NewDecoder(rw)
	if err := decoder.Decode(&t); err!=nil {
		log.Println(err)
	}

	// Read from the connection until EOF. Expect a command name as the
	// next input. Call the handler that is registered for this command.
	//for {
	//	log.Print("Receive command '")
	//	cmd, err := rw.ReadString('\n')
	//	switch {
	//	case err == io.EOF:
	//		log.Println("Reached EOF - close this connection.\n   ---")
	//		return
	//	case err != nil:
	//		log.Println("\nError reading command. Got: '"+cmd+"'\n", err)
	//		return
	//	}
	//	// Trim the request string - ReadString does not strip any newlines.
	//	cmd = strings.Trim(cmd, "\n ")
	//	log.Println(cmd + "'")
	//
	//	// Fetch the appropriate handler function from the 'handler' map and call it.
	//	e.m.RLock()
	//	handleCommand, ok := e.handler[cmd]
	//	e.m.RUnlock()
	//	if !ok {
	//		log.Println("Command '" + cmd + "' is not registered.")
	//		return
	//	}
	//	handleCommand(rw)
	//}
}


/* Now let's create two handler functions. The easiest case is where our
ad-hoc protocol only sends string data.

The second handler receives and processes a struct that was send as GOB data.
*/

// handleStrings handles the "STRING" request.
func handleStrings(rw *bufio.ReadWriter) {
	// Receive a string.
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

// handleGob handles the "GOB" request. It decodes the received GOB data
// into a struct.
func handleGob(rw *bufio.ReadWriter) {
	log.Print("Receive GOB data:")
	var data complexData
	// Create a decoder that decodes directly into a struct variable.
	dec := gob.NewDecoder(rw)
	err := dec.Decode(&data)
	if err != nil {
		log.Println("Error decoding GOB data:", err)
		return
	}
	// Print the complexData struct and the nested one, too, to prove
	// that both travelled across the wire.
	log.Printf("Outer complexData struct: \n%#v\n", data)
	log.Printf("Inner complexData struct: \n%#v\n", data.C)
}



// server listens for incoming requests and dispatches them to
// registered handler functions.
func server() error {
	endpoint := NewEndpoint()

	// Add the handle funcs.
	endpoint.AddHandleFunc("STRING", handleStrings)
	endpoint.AddHandleFunc("GOB", handleGob)

	// Start listening.
	return endpoint.Listen()
}
