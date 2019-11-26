package utils

import (
	"encoding/json"

	"github.com/daominah/gomicrokit/log"
	"github.com/daominah/socketcluster-client-go/scclient/models"
	"github.com/shamaton/msgpack"
)

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
		result, err = msgpack.Encode(newV)

	case models.EmitEvent:
		channel, ok := v.Data.(models.Channel)
		if !ok {
			result, err = msgpack.Encode(v)
			break
		}
		array := []interface{}{channel.Channel, channel.Data}
		if v.Cid != 0 {
			array = append(array, v.Cid)
		}
		var newV interface{}
		switch v.Event {
		case "#publish":
			newV = ModifiedMsgP{P: array}
		default:
			newV = ModifiedMsgE{E: array}
		}
		result, err = msgpack.Encode(newV)

	case models.ReceiveEvent:
		newV := ModifiedMsgR{R: []interface{}{v.Rid, v.Error, v.Data}}
		result, err = msgpack.Encode(newV)

	default:
		result, err = msgpack.Encode(v)
	}
	//log.Debugf("msgpacker Marshal out %v: %v\n", err, result)
	return result, err
}

func (m Msgpacker) Unmarshal(data []byte, i interface{}) error {
	//log.Debugf("msgpacker Unmarshal in: %#v", i)
	return msgpack.Decode(data, i)
}
