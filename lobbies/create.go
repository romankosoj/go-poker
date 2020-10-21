package lobbies

import (
	"log"

	"github.com/JohnnyS318/go-poker/lobby"
)

func (l *LobbyManager) CreateNew(c chan int) int {

	le := len(l.Lobbies)
	if le >= l.MaxCount {
		return -1
	}

	lobby := lobby.NewLobby()
	lobby.Callback = func() {
		c <- le
	}

	// lobbies can still be created
	l.Lobbies = append(l.Lobbies, *lobby)

	log.Printf("Created new lobby %v", le)

	return le
}
