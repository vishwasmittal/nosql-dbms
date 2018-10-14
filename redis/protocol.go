package redis

/*
Attributes of a protocol
	- incoming message handler
	- output response serializer
*/

//var COMMANDS map[string]string = {
//"get": "GET",
//"set": "SET",
//"del": "DEL",
//"evict": "EVICT"
//}

type RequestProtocol struct {
	Command string

	// data will either be a string (in case of GET and DEL) or a map (in case of SET)
	Data string
}

type ResponseProtocol struct {
	Error string
	Data  string
}
