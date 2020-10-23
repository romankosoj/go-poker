package handlers

import (
	"log"
	"net/http"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Lobby) Join(rw http.ResponseWriter, r *http.Request) {

	log.Printf("/join called")

	conn, err := upgrader.Upgrade(rw, r, nil)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadGateway)
	}

	playerConn := NewPlayerConn(conn)
	playerConn.OnClose = func(err error) {
		log.Printf("Closing err %v", err)
		return
	}

	go playerConn.reader()

	go playerConn.writer()

	log.Printf("Waiting for join Event")

	raw := <-playerConn.In

	joinEvent, err := events.ToJoinEvent(raw)

	if err != nil {
		log.Printf("joinEvent was invalid %v", err)
		playerConn.Out <- models.NewEvent("VALIDATION_FAILED", "The joining event was not as the server expected").ToRaw()
		conn.Close()
		http.Error(rw, "VALIDATION_FAILED. The joining event was not as the server expected", http.StatusBadRequest)
		return
	}

	player := models.NewPlayer(joinEvent.Username, joinEvent.ID, joinEvent.BuyIn, playerConn.In, playerConn.Out)

	lobby, _, err := h.Lobbies.ManagePlayer(player)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadGateway)
		return
	}

	playerConn.OnClose = func(err error) {
		log.Printf("Removed player from lobby")
		lobby.RemovePlayerByID(joinEvent.ID)
	}

}
