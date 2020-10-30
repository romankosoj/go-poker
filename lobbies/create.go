package lobbies

import (
	"errors"
	"log"

	"github.com/JohnnyS318/go-poker/lobby"
)

func (l *LobbyManager) CreateNew(c chan string) (string, error) {

	le := len(l.Lobbies)
	if le >= l.MaxCount {
		return "", errors.New("Maximum Lobby count already passed")
	}

	lobby := lobby.NewLobby(c)

	l.Lobbies[lobby.LobbyID] = lobby
	log.Printf("Created new lobby [%v]", lobby.LobbyID)

	return lobby.LobbyID, nil
}
