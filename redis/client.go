package redis

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net"
)

/*
Attributes of a client
	- host
	- port (default: 16379
	- provide functions for each of the server supported commands
*/
type Client struct {
}

// Open connects to a TCP Address.
// It returns a TCP connection armed with a timeout and wrapped into a
// buffered ReadWriter.
func Open(addr string) (*bufio.ReadWriter, error) {
//func Open(addr string) (net.Conn, error) {
	// Dial the remote process.
	// Note that the local port is chosen on the fly. If the local port
	// must be a specific one, use DialTCP() instead.
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn)), nil
	//return conn, nil
}

/*
## The client and server functions

With all this in place, we can now set up client and server functions.

The client function connects to the server and sends STRING and GOB requests.

The server starts listening for requests and triggers the appropriate handlers.
*/

// client is called if the app is called with -connect=`ip addr`.
func ClientFunc(ip string) error {

	//var t = temp_interface{"yolo", 1020}
	//log.Println(t)
	rw, err := Open(ip + Port)
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip+Port)
	}

	//encoder:=json.NewEncoder(rw)
	//encoder.Encode(t)
	//rw.Flush()
	//var rec_t temp_interface
	//decoder:=json.NewDecoder(rw)
	//decoder.Decode(&rec_t)
	//
	//log.Println(rec_t)

	rw.WriteString("This is the client string\n")
	rw.Flush()
	data, err := rw.ReadString('\n')
	fmt.Println(data)
	return nil

	//// Some test data. Note how GOB even handles maps, slices, and
	//// recursive data structures without problems.
	//testStruct := complexData{
	//	N: 23,
	//	S: "string data",
	//	M: map[string]int{"one": 1, "two": 2, "three": 3},
	//	P: []byte("abc"),
	//	C: &complexData{
	//		N: 256,
	//		S: "Recursive structs? Piece of cake!",
	//		M: map[string]int{"01": 1, "10": 2, "11": 3},
	//	},
	//}
	//
	//// Open a connection to the server.
	//rw, err := Open(ip + Port)
	//if err != nil {
	//	return errors.Wrap(err, "Client: Failed to open connection to "+ip+Port)
	//}
	//
	//// Send a STRING request.
	//// Send the request name.
	//// Send the data.
	//log.Println("Send the string request.")
	//n, err := rw.WriteString("STRING\n")
	//if err != nil {
	//	return errors.Wrap(err, "Could not send the STRING request ("+strconv.Itoa(n)+" bytes written)")
	//}
	//n, err = rw.WriteString("Additional data.\n")
	//if err != nil {
	//	return errors.Wrap(err, "Could not send additional STRING data ("+strconv.Itoa(n)+" bytes written)")
	//}
	//log.Println("Flush the buffer.")
	//err = rw.Flush()
	//if err != nil {
	//	return errors.Wrap(err, "Flush failed.")
	//}
	//
	//// Read the reply.
	//log.Println("Read the reply.")
	//response, err := rw.ReadString('\n')
	//if err != nil {
	//	return errors.Wrap(err, "Client: Failed to read the reply: '"+response+"'")
	//}
	//
	//log.Println("STRING request: got a response:", response)
	//
	//// Send a GOB request.
	//// Create an encoder that directly transmits to `rw`.
	//// Send the request name.
	//// Send the GOB.
	//log.Println("Send a struct as GOB:")
	//log.Printf("Outer complexData struct: \n%#v\n", testStruct)
	//log.Printf("Inner complexData struct: \n%#v\n", testStruct.C)
	//enc := gob.NewEncoder(rw)
	//n, err = rw.WriteString("GOB\n")
	//if err != nil {
	//	return errors.Wrap(err, "Could not write GOB data ("+strconv.Itoa(n)+" bytes written)")
	//}
	//err = enc.Encode(testStruct)
	//if err != nil {
	//	return errors.Wrapf(err, "Encode failed for struct: %#v", testStruct)
	//}
	//err = rw.Flush()
	//if err != nil {
	//	return errors.Wrap(err, "Flush failed.")
	//}
	//return nil
}
