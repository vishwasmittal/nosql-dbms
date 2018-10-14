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

//SupportedDataTypes := map[string]string {}
//"float32": "float32"
//"int": "int"
//"string": "string"
//}

type DynamicDataStruct struct {
	DType string
	Value interface{}
}

type DataStruct struct {
	Key  string
	Data DynamicDataStruct
}

var EmptyDynamicDataSctuct = DynamicDataStruct{"", nil}
var EmptyDataStruct DataStruct = DataStruct{"", EmptyDynamicDataSctuct}

func NewDataStruct(key string, dType string, value interface{}) DataStruct {
	return DataStruct{
		key, DynamicDataStruct{dType, value}}
}

type RequestProtocol struct {
	Command string
	Data    DataStruct
}

type ResponseProtocol struct {
	Error string
	Data  DataStruct
}
