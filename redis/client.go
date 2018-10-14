package redis

import (
	"encoding/json"
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

//func Open(addr string) (*bufio.ReadWriter, error) {
func Open(addr string) (net.Conn, error) {
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return conn, nil
}

func ClientFunc(ip string) error {
	rw, err := Open(ip + Port)
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ip+Port)
	}

	var server_c = RequestProtocol{"GET", "Hello"}
	encoder := json.NewEncoder(rw)
	decoder := json.NewDecoder(rw)

	fmt.Println(server_c)
	e := encoder.Encode(server_c)
	if e != nil {
		fmt.Println("Error occured here!!!")
	}

	var recv_c ResponseProtocol
	decoder.Decode(&recv_c)
	log.Println(recv_c)
	return nil
}
