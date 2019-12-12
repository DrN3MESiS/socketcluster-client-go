package scclient

import (
	"reflect"
	"sync"

	logging "github.com/sacOO7/go-logger"
)

type Empty struct{}

var scLogger = logging.GetLogger(reflect.TypeOf(Empty{}).PkgPath()).SetLevel(logging.OFF)

type Listener struct {
	emitAckListener map[int][]interface{}
	onListener      map[string]func(eventName string, data interface{})
	onAckListener   map[string]func(eventName string, data interface{}, ack func(error interface{}, data interface{}))
	// u fucking noob Go coder sacOO7
	mutex           sync.Mutex
}

func Init() Listener {
	return Listener{
		emitAckListener: make(map[int][]interface{}),
		onListener:      make(map[string]func(eventName string, data interface{})),
		onAckListener:   make(map[string]func(eventName string, data interface{}, ack func(error interface{}, data interface{}))),
	}
}

func (listener *Listener) putEmitAck(id int, eventName string, ack func(eventName string, error interface{}, data interface{})) {
	listener.mutex.Lock()
	listener.emitAckListener[id] = []interface{}{eventName, ack}
	listener.mutex.Unlock()
}

func (listener *Listener) handleEmitAck(id int, error interface{}, data interface{}) {
	listener.mutex.Lock()
	ackObject := listener.emitAckListener[id]
	delete(listener.emitAckListener, id)
	listener.mutex.Unlock()
	if ackObject != nil {
		eventName := ackObject[0].(string)
		scLogger.Trace.Println("Ack received for event :: ", eventName)
		ack := ackObject[1].(func(eventName string, error interface{}, data interface{}))
		ack(eventName, error, data)
	} else {
		scLogger.Warning.Println("Ack function not found for rid :: ", id)
	}
}

func (listener *Listener) putOnListener(eventName string, onListener func(eventName string, data interface{})) {
	listener.mutex.Lock()
	listener.onListener[eventName] = onListener
	listener.mutex.Unlock()
}

func (listener *Listener) handleOnListener(eventName string, data interface{}) {
	listener.mutex.Lock()
	on := listener.onListener[eventName]
	listener.mutex.Unlock()
	if on != nil {
		on(eventName, data)
	}
}

func (listener *Listener) putOnAckListener(eventName string, onAckListener func(eventName string, data interface{}, ack func(error interface{}, data interface{}))) {
	listener.mutex.Lock()
	listener.onAckListener[eventName] = onAckListener
	listener.mutex.Unlock()
}

func (listener *Listener) handleOnAckListener(eventName string, data interface{}, ack func(error interface{}, data interface{})) {
	listener.mutex.Lock()
	onAck := listener.onAckListener[eventName]
	listener.mutex.Unlock()
	if onAck != nil {
		onAck(eventName, data, ack)
	}
}

func (listener *Listener) hasEventAck(eventName string) bool {
	listener.mutex.Lock()
	temp := listener.onAckListener[eventName]
	listener.mutex.Unlock()
	return temp != nil
}
