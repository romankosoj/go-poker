package main

import (
	"log"
	"net/http"

	"github.com/JohnnyS318/go-poker/handlers"
	"github.com/JohnnyS318/go-poker/lobbies"
	"github.com/gorilla/mux"
)

func main() {

	// Setup
	r := mux.NewRouter()

	// Setup Lobby map

	lobbyManager := lobbies.NewManager(10)

	// Setup Handlers
	playerHandler := handlers.NewLobbyHandler(lobbyManager)

	// each lobby creates a new game and waites for at least 3 players

	// /join {PlayerId, PlayerUsername} searches through the lobby map or creates one if none are found and the maximum lobby count is not succeeded.
	r.HandleFunc("/join", playerHandler.Join)

	// join upgrades connection to websocket

	// game starts and communicates only via the websocket connection

	log.Printf("Serve on Port %v", 8080)

	log.Printf(http.ListenAndServe(":8080", r).Error())

}
