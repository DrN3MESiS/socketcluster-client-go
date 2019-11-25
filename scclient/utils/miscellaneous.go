package utils

import (
	"encoding/json"
	"fmt"

	"github.com/rgamba/evtwebsocket"
	"github.com/vmihailenco/msgpack"
)

var marshaller Marshaller
var unmarshaller Unmarshaller

func init() {
	marshaller, unmarshaller = Jsoner{}, Jsoner{}
	//marshaller, unmarshaller = Msgpacker{}, Msgpacker{}
}

type Marshaller interface {
	Marshal(interface{}) ([]byte, error)
}

type Unmarshaller interface {
	Unmarshal(data []byte, v interface{}) error
}

type Jsoner struct{}

func (m Jsoner) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (m Jsoner) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type Msgpacker struct{}

func (m Msgpacker) Marshal(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)
}

func (m Msgpacker) Unmarshal(data []byte, v interface{}) error {
	return msgpack.Unmarshal(data, v)
}

func PrintMessage(message string) {
	fmt.Println(message)
}

func IsEqual(s string, b []byte) bool {
	if len(s) != len(b) {
		return false
	}
	for i, x := range b {
		if x != s[i] {
			return false
		}
	}
	return true
}

func CreateMessageFromString(message string) evtwebsocket.Msg {
	return evtwebsocket.Msg{
		Body: []byte(message),
	}
}

func CreateMessageFromByte(message []byte) evtwebsocket.Msg {
	return evtwebsocket.Msg{
		Body: message,
	}
}

func SerializeData(data interface{}) []byte {
	b, _ := marshaller.Marshal(data)
	return b
}

func SerializeDataIntoString(data interface{}) string {
	b, _ := marshaller.Marshal(data)
	return string(b)
}

func DeserializeData(data []byte) (jsonObject interface{}) {
	unmarshaller.Unmarshal(data, &jsonObject)
	return
}

func DeserializeDataFromString(data string) (jsonObject interface{}) {
	unmarshaller.Unmarshal([]byte(data), &jsonObject)
	return
}
