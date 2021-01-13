package websocket

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"ptop/internal/websocket"
	"time"
)

func Setup(server *websocket.Server) error {
	router := mux.NewRouter()
	handler := &Handler{Server: server}

	router.HandleFunc("/new", handler.NewWS).Methods("GET")
	router.HandleFunc("/peers", handler.GetPeers)

	log.Println("Listening on", server.Port)
	s := &http.Server{
		Addr:           ":" + server.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := s.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
