package parser

import (
	"fmt"

	"strconv"

	"github.com/daominah/gomicrokit/log"
)

func Parse(rid int, cid int, event interface{}) MessageType {
	if event != nil {
		if event == "#publish" {
			return PUBLISH

		} else if event == "#removeAuthToken" {
			return REMOVETOKEN

		} else if event == "#setAuthToken" {
			return SETTOKEN

		} else {
			return EVENT
		}
	} else if rid == 1 {
		return ISAUTHENTICATED

	} else {
		return ACKRECEIVE
	}
}

func GetMessageDetails(message interface{}) (data interface{}, rid int, cid int, eventname interface{}, error interface{}) {
	//Converting given message into map, with keys and values to that we can parse it

	itemsMap, ok := message.(map[string]interface{})
	if !ok {
		log.Infof("unexpected message type: %T", message)
		return
	}

	for k, v := range itemsMap {
		switch k {
		case "data":
			data = v
		case "rid":
			rStr := fmt.Sprintf("%v", v)
			ridInt64, _ := strconv.ParseInt(rStr, 10, 64)
			rid = int(ridInt64)
		case "cid":
			rStr := fmt.Sprintf("%v", v)
			ridInt64, _ := strconv.ParseInt(rStr, 10, 64)
			cid = int(ridInt64)
		case "event":
			eventname = v
		case "error":
			error = v
		}
	}

	return
}
