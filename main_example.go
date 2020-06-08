package main

import (
	"log"
	"os"
	"time"

	"github.com/daominah/gomicrokit/websocket"
	"github.com/daominah/socketcluster-client-go/scclient"
)

var channelName0 = "channel0"

func print3args(ev string, err interface{}, data interface{}) {
	log.Println("print1", ev, err, data)
}
func print2args(ev string, data interface{}) {
	log.Println("print2", ev, data)
}

func onConnect(client *scclient.Client) {
	log.Println("onConnect")
	client.PublishAck(channelName0, map[string]interface{}{"fuck": 1}, print3args)
	time.Sleep(1 * time.Second)
	client.SubscribeAck(channelName0, print3args)
}

func onDisconnect(client *scclient.Client, err error) {
	log.Println("onDisconnect:", err)
	os.Exit(1)
}

func onConnectError(client *scclient.Client, err error) {
	log.Println("onConnectError:", err)
	os.Exit(1)
}

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	websocket.LOG = true
	client := scclient.New("ws://10.100.50.100:8000/socketcluster/")
	client.SetBasicListener(onConnect, onConnectError, onDisconnect)
	client.EnableLogging()
	client.OnChannel(channelName0, print2args)
	go client.Connect()
	select {}
}
