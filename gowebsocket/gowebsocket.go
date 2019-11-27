package gowebsocket

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/daominah/gomicrokit/websocket"
	goraws "github.com/gorilla/websocket"
)

func init() {
	websocket.LOG = false
	websocket.SetWebsocketConfig(
		20*time.Second, 20*time.Second, 8*time.Second, 65536)
}

type Socket struct {
	Conn              *websocket.Connection
	WebsocketDialer   *goraws.Dialer
	Url               string
	ConnectionOptions ConnectionOptions
	RequestHeader     http.Header
	OnConnected       func(socket Socket)
	OnTextMessage     func(message string, socket Socket)
	OnBinaryMessage   func(data []byte, socket Socket)
	OnConnectError    func(err error, socket Socket)
	OnDisconnected    func(err error, socket Socket)
	OnPingReceived    func(data string, socket Socket)
	OnPongReceived    func(data string, socket Socket)
	IsConnected       bool
	sendMu            *sync.Mutex // Prevent "concurrent write to websocket connection"
	receiveMu         *sync.Mutex
}

type ConnectionOptions struct {
	UseCompression bool
	UseSSL         bool
	Proxy          func(*http.Request) (*url.URL, error)
	Subprotocols   []string
}

func New(url string) Socket {
	return Socket{
		Url:           url,
		RequestHeader: http.Header{},
		ConnectionOptions: ConnectionOptions{
			UseCompression: false,
			UseSSL:         true,
		},
		WebsocketDialer: &goraws.Dialer{},
		sendMu:          &sync.Mutex{},
		receiveMu:       &sync.Mutex{},
	}
}

func (socket *Socket) setConnectionOptions() {
	socket.WebsocketDialer.EnableCompression = socket.ConnectionOptions.UseCompression
	socket.WebsocketDialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: socket.ConnectionOptions.UseSSL}
	socket.WebsocketDialer.Proxy = socket.ConnectionOptions.Proxy
	socket.WebsocketDialer.Subprotocols = socket.ConnectionOptions.Subprotocols
}

type onReadHeandler struct{ socket *Socket }

func (h onReadHeandler) Handle(cid websocket.ConnectionId, msgType int, msg []byte) {
	if msgType == goraws.BinaryMessage {
		if h.socket.OnBinaryMessage != nil {
			h.socket.OnBinaryMessage(msg, *h.socket)
		}
	} else if msgType == goraws.TextMessage {
		if h.socket.OnTextMessage != nil {
			h.socket.OnTextMessage(string(msg), *h.socket)
		}
	} else {
		// pass
	}
}

func (socket *Socket) Connect() {
	socket.setConnectionOptions()

	goraConn, _, err := socket.WebsocketDialer.Dial(socket.Url, socket.RequestHeader)
	if err != nil {
		socket.IsConnected = false
		if socket.OnConnectError != nil {
			socket.OnConnectError(err, *socket)
		}
		return
	}

	socket.Conn = websocket.NewConnection(goraConn, onReadHeandler{socket: socket})

	if socket.OnConnected != nil {
		socket.IsConnected = true
		socket.OnConnected(*socket)
	}

	if socket.OnDisconnected != nil {
		go func() {
			<-socket.Conn.ClosedChan
			socket.IsConnected = false
			socket.OnDisconnected(err, *socket)
		}()
	}
}

func (socket *Socket) SendText(message string) {
	go socket.Conn.Write(message)
}

func (socket *Socket) SendBinary(message []byte) {
	go socket.Conn.WriteBytes(message)
}

func (socket *Socket) Close() {
	socket.Conn.Close()
}
