package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Peer struct {
	Socket *websocket.Conn
	Send   chan []byte
	Target  string
}

func NewPeer(conn *websocket.Conn, target string) *Peer {
	return &Peer{
		Socket: conn,
		Send:   make(chan []byte),
		Target: target,
	}
}

type Message struct {
	Event   string `json:"event"`
	Content string `json:"content"`
}

func (peer *Peer) Read(server *Server) {
	defer func() {
		_ = peer.Socket.Close()
	}()

	log.Println("Listening for a message")

	for {
		_, message, err := peer.Socket.ReadMessage()
		log.Println("Message handled")
		if err != nil {
			_ = peer.Socket.Close()
			break
		}

		m := Message{}
		err = json.Unmarshal(message, &m)
		if err != nil {
			fmt.Println("Peer Read() err : " + err.Error())
			continue
		}

		if m.Event == "new_addr" {
			server.ConnectToAddress(m.Content, true)
		}
	}
}

func (peer *Peer) Write() {
	defer func() {
		_ = peer.Socket.Close()
	}()

	for {
		select {
		case message, ok := <- peer.Send:
			if !ok {
				_ = peer.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			_ = peer.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
