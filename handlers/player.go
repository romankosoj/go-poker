package handlers

import "github.com/JohnnyS318/go-poker/lobbies"

type Lobby struct {
	Lobbies *lobbies.LobbyManager
}

func NewLobbyHandler(lobbyManager *lobbies.LobbyManager) *Lobby {
	return &Lobby{
		Lobbies: lobbyManager,
	}
}
