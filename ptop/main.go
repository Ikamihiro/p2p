package main

import (
	"flag"
	"fmt"
	"log"
	"ptop/internal/websocket"
	websocket2 "ptop/pkg/websocket"
)

var (
	port string
	seed string
)

func initFlags() {
	flag.StringVar(&port, "p", "8080", "listen port")
	flag.StringVar(&seed, "seed", "127.0.0.1:8080", "seed ip:port")
	flag.Parse()
}

func main() {
	fmt.Println("My first Peer To Peer project")

	initFlags()

	server := websocket.NewServer(port)

	if ("127.0.0.1:" + port) != seed {
		go server.ConnectToAddress(seed, false)
	}

	log.Fatal(websocket2.Setup(server))
}
