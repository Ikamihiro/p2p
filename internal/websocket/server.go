package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
)

type Server struct {
	Peers []*Peer
	Port string
}

func NewServer(port string) *Server {
	return &Server{
		Peers: []*Peer{},
		Port:  port,
	}
}

func (server *Server) ConnectToAddress(address string, isBroadcast bool) {
	rawQ := "port=" + server.Port

	if isBroadcast {
		rawQ += ";brdcst=1"
	}

	u := url.URL{
		Scheme: "ws",
		Host: address,
		Path: "/new",
		RawQuery: rawQ,
	}

	log.Printf("Connect with %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("ConnectionToAddr err: " + err.Error())
		return
	}

	newPeer := server.AppendNewPeer(conn, address)
	go newPeer.Read(server)
	go newPeer.Write()
}

func (server *Server) BroadcastAddress(target string) {
	m :=  &Message{
		Event:   "new_addr",
		Content: target,
	}
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println("BroadcastAddr json error : " + err.Error())
		return
	}

	for _, p := range server.Peers {
		p.Send <- b
	}
}

func (server *Server) AppendNewPeer(conn *websocket.Conn, target string) *Peer {
	log.Printf("Adding Peer %s", target)
	peer := NewPeer(conn, target)
	server.Peers = append(server.Peers, peer)
	return peer
}
