package utils

import (
	"encoding/json"
	"errors"

	"github.com/daominah/gomicrokit/log"
	"github.com/daominah/socketcluster-client-go/scclient/models"
	"github.com/shamaton/msgpack"
)

func msgpackEncode(v interface{}) ([]byte, error) {
	return msgpack.Encode(v)
	//return msgpack.Marshal(v)
}

func msgpackDecode(data []byte, v interface{}) error {
	return msgpack.Decode(data, v)
	//return msgpack.Unmarshal(data, v)
}

type Jsoner struct{}

func (m Jsoner) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (m Jsoner) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type Msgpacker struct{}

type ModifiedMsgP struct {
	P []interface{} `json:"p" msgpack:"p"`
}
type ModifiedMsgE struct {
	E []interface{} `json:"e" msgpack:"e"`
}
type ModifiedMsgR struct {
	R []interface{} `json:"r" msgpack:"r"`
}

func (m Msgpacker) Marshal(i interface{}) (result []byte, err error) {
	_ = log.Debugf
	//log.Debugf("msgpacker Marshal in: %#v", i)
	switch v := i.(type) {
	case models.HandShake:
		newV := ModifiedMsgE{E: []interface{}{v.Event, v.Data, v.Cid}}
		result, err = msgpackEncode(newV)

	case models.EmitEvent:
		channel, ok := v.Data.(models.Channel)
		if !ok {
			result, err = msgpackEncode(v)
			break
		}
		var newV interface{}
		switch v.Event {
		case "#publish":
			array := []interface{}{channel.Channel, channel.Data}
			if v.Cid != 0 {
				array = append(array, v.Cid)
			}
			newV = ModifiedMsgP{P: array}
		default:
			array := []interface{}{v.Event, channel}
			if v.Cid != 0 {
				array = append(array, v.Cid)
			}
			//log.Println("a", array)
			newV = ModifiedMsgE{E: array}
		}
		result, err = msgpackEncode(newV)

	case models.ReceiveEvent:
		newV := ModifiedMsgR{R: []interface{}{v.Rid, v.Error, v.Data}}
		result, err = msgpackEncode(newV)

	default:
		result, err = msgpackEncode(v)
	}
	//log.Debugf("msgpacker Marshal out %v: %v, %s \n", err, result, result)
	return result, err
}

func (m Msgpacker) Unmarshal(data []byte, i interface{}) error {
	err := msgpackDecode(data, i)
	if err != nil {
		return err
	}
	i1, ok := i.(*interface{})
	if !ok {
		return errors.New("un expected type i")
	}
	//log.Debugf("msgpacker Unmarshal in: %#v", *i1)
	i2, ok := (*i1).(map[interface{}]interface{})
	if !ok {
		return errors.New("un expected type *i1")
	}
	i3 := make(map[string]interface{})
	for k, v := range i2 {
		ks, ok := k.(string)
		if !ok {
			continue
		}
		i3[ks] = v
	}
	var arrayI interface{}
	var field string
	ok = true
	for _, ifield := range []string{"p", "e", "r"} {
		arrayI, ok = i3[ifield]
		if ok {
			field = ifield
			break
		}
	}
	if !ok {
		return errors.New("field should be p, e or r")
	}
	array, ok := arrayI.([]interface{})
	if !ok {
		return errors.New("un expected type arrayI")
	}
	i4 := make(map[string]interface{})
	if field == "r" {
		if len(array) != 3 {
			return errors.New("un expected len array")
		}
		i4["rid"] = array[0]
		i4["error"] = array[1]
		i4["data"] = array[2]
	} else { //if field == "p" || field == "e"
		if len(array) < 2 {
			return errors.New("un expected len array <2")
		}
		i4["channel"] = array[0]
		i4["data"] = array[1]
		if len(array) == 3 {
			i4["cid"] = array[2]
		}
	}
	i4DataI, ok := i4["data"].(map[interface{}]interface{})
	if ok {
		i4Data := make(map[string]interface{})
		for k, v := range i4DataI {
			ks, ok := k.(string)
			if !ok {
				continue
			}
			i4Data[ks] = v
		}
		i4["data"] = i4Data
	}
	*i1 = i4
	//log.Debugf("msgpacker Unmarshal out: %#v", *i1)
	return nil
}
