package redis

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net"
)


func Open(addr string) (net.Conn, error) {
	log.Println("Dial " + addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, errors.Wrap(err, "Dialing "+addr+" failed")
	}
	return conn, nil
}

var ipG string

func testCommand(command string, dataStruct DataStruct) error {
	conn, err := Open(ipG + Port)
	if err != nil {
		return errors.Wrap(err, "Client: Failed to open connection to "+ipG+Port)
	}

	var server_c = RequestProtocol{command, dataStruct}
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)
	fmt.Println("Request data: ", server_c)
	e := encoder.Encode(server_c)
	if e != nil {
		fmt.Println("Error occured here: 2")
	}
	var recv_c ResponseProtocol
	decoder.Decode(&recv_c)
	fmt.Println("Response Data: ", recv_c)

	return nil
}

func ClientFunc(ip string) error {

	testCommand("GET", NewDataStruct("1", "", nil))
	testCommand("SET", NewDataStruct("1", "string", "111111111111"))
	testCommand("GET", NewDataStruct("1", "", nil))
	testCommand("SET", NewDataStruct("2", "string", "22222222222"))
	testCommand("GET", NewDataStruct("2", "", nil))
	testCommand("DEL", NewDataStruct("1", "", nil))
	testCommand("GET", NewDataStruct("1", "", nil))
	testCommand("EVICT", EmptyDataStruct)
	testCommand("GET", NewDataStruct("2", "", nil))
	return nil
}
