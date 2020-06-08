package models

type EmitEvent struct {
	Event string      `json:"event" msgpack:"data"`
	Data  interface{} `json:"data" msgpack:"data"`
	Cid   int         `json:"cid" msgpack:"data"`
}

type ReceiveEvent struct {
	Data  interface{} `json:"data" msgpack:"data"`
	Error interface{} `json:"error" msgpack:"error"`
	Rid   int         `json:"rid" msgpack:"rid" `
}

type Channel struct {
	Channel string      `json:"channel" msgpack:"channel"`
	Data    interface{} `json:"data,omitempty" msgpack:"data"`
}

func GetEmitEventObject(eventname string, data interface{}, messageId int) EmitEvent {
	return EmitEvent{
		Event: eventname,
		Data:  data,
		Cid:   messageId,
	}
}

func GetReceiveEventObject(data interface{}, error interface{}, messageId int) ReceiveEvent {
	return ReceiveEvent{
		Data:  data,
		Error: error,
		Rid:   messageId,
	}
}

func GetChannelObject(data interface{}) Channel {
	channelObject := data.(map[string]interface{})
	return Channel{Channel: channelObject["channel"].(string), Data: channelObject["data"]}
}

func GetSubscribeEventObject(channelName string, messageId int) EmitEvent {
	return EmitEvent{
		Event: "#subscribe",
		Data:  Channel{Channel: channelName},
		Cid:   messageId,
	}
}

func GetUnsubscribeEventObject(channelName string, messageId int) EmitEvent {
	return EmitEvent{
		Event: "#unsubscribe",
		Data:  Channel{Channel: channelName},
		Cid:   messageId,
	}
}

func GetPublishEventObject(channelName string, data interface{}, messageId int) EmitEvent {
	return EmitEvent{
		Event: "#publish",
		Data:  Channel{Channel: channelName, Data: data},
		Cid:   messageId,
	}
}
