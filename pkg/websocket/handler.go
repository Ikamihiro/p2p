package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	internal "ptop/internal/websocket"
	"strings"
)

type Handler struct {
	Server *internal.Server
}

func (handler *Handler) GetPeers(res http.ResponseWriter, req *http.Request) {
	var addresses []string
	for _, p := range handler.Server.Peers {
		addresses = append(addresses, p.Target)
	}

	b, err := json.Marshal(addresses)
	if err != nil {
		fmt.Println(err.Error())
	}

	_, _ = res.Write(b)
}

func (handler *Handler) NewWS(res http.ResponseWriter, req *http.Request) {
	q := req.URL.Query()
	rPort, ok := q["port"]

	log.Printf("Endpoint NewWS: query: %s", q.Encode())

	if !ok {
		fmt.Println("url value port is nil")
		http.NotFound(res, req)
		return
	}

	ip := strings.Split(req.RemoteAddr, ":")
	target := ip[0] + ":" + rPort[0]

	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		fmt.Println("new client error: " + err.Error())
		http.NotFound(res, req)
		return
	}

	v, ok := q["brdcst"]
	if !ok {
		if len(v) == 0 {
			handler.Server.BroadcastAddress(target)
		}
	}

	newPeer := handler.Server.AppendNewPeer(conn, target)
	go newPeer.Write()
	go newPeer.Read(handler.Server)
}
