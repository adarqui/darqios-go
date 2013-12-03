package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
	"strings"
	"encoding/json"
)

func (M *Main) WS_Init() {
	DebugLn("WS_Init:Initialized")
}

type BMSG_DATA struct {
	Host string
	Channel string
	Data interface{}
}

type BMSG struct {
	Channel string
	Data []byte
}

/* WSOP = WebSocketOP */
const (
	WSOP_REQ_REGISTER = "register"
	WSOP_REQ_UNREGISTER = "unregister"
)

type WSOP struct {
	Op string
	Channel string
	Data interface{}
}

// WS maintains the set of active connections and broadcasts messages to the
// connections.
type WS struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
//	broadcast chan []byte
	broadcast chan BMSG

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var HUB = WS{
	//⇥ broadcast:   make(chan []byte),
	broadcast: make(chan BMSG),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}


func (HUB *WS) WS_Run() {
	for {
		select {
		case c := <-HUB.register:
			HUB.connections[c] = true
		case c := <-HUB.unregister:
			delete(HUB.connections, c)
			close(c.send)
		case m := <-HUB.broadcast:
			for c := range HUB.connections {
				log.Printf("BROADCAST: m.channel: %s, channels: %q\n", m.Channel, c.channels)
				if truth, ok := c.channels[m.Channel]; ok {
					if truth != true {
						continue
					}
					select {
					case c.send <- m.Data:
					default:
						close(c.send)
						delete(HUB.connections, c)
					}
				}
			}
		}
	}
}


const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// connection is an middleman between the websocket connection and the WS.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	channels map[string]bool
}

// readPump pumps messages from the websocket connection to the WS.
func (c *connection) WS_readPump() {
	defer func() {
		HUB.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		log.Printf("READPUMP: %v\n", message)
		bmsg := BMSG{}
		bmsg.Channel = "cli"
		bmsg.Data = message

		var wsop WSOP
		dec := json.NewDecoder(strings.NewReader(string(message)))
		dec.Decode(&wsop)

		log.Printf("WSOP: %q %q\n", wsop, message)
//		log.Fatal(op)

		var truth bool


		/*
		 * WEBSOCKETO OPERATIONS
		 * register
		 * unregister
		 */
		switch wsop.Op {
			case WSOP_REQ_REGISTER: truth = true
			case WSOP_REQ_UNREGISTER: truth = false
			default: return
		}

		if wsop.Channel != "" {
			c.channels[wsop.Channel] = truth
		}

//		HUB.broadcast <- bmsg
	}
}

// write writes a message with the given message type and payload.
func (c *connection) WS_write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the WS to the websocket connection.
func (c *connection) WS_writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.WS_write(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("TEXTMESSAGE: %v\n", message)
			if err := c.WS_write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			log.Printf("PINGMESSAGE!\n")
			if err := c.WS_write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serverWs handles webocket requests from the peer.
func WS_Serve(M *Main, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws, channels: make(map[string]bool)}
	HUB.register <- c
	go c.WS_writePump()
	c.WS_readPump()
}
