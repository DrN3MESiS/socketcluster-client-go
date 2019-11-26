package utils

import (
	"fmt"

	"github.com/rgamba/evtwebsocket"
)

var marshaller Marshaller
var unmarshaller Unmarshaller

func init() {
	//SetMarshaller(false)
	SetMarshaller(true)
}

func SetMarshaller(isCodecMinBin bool) {
	if isCodecMinBin {
		marshaller, unmarshaller = Msgpacker{}, Msgpacker{}
	} else {
		marshaller, unmarshaller = Jsoner{}, Jsoner{}
	}
}

type Marshaller interface {
	Marshal(interface{}) ([]byte, error)
}

type Unmarshaller interface {
	Unmarshal(data []byte, v interface{}) error
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
