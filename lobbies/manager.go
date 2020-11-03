package lobbies

import (
	"errors"
	"log"
	"math/rand"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/lobby"
	"github.com/JohnnyS318/go-poker/models"
)

type LobbyManager struct {
	Lobbies        map[string]*lobby.Lobby
	Capacity       map[string]int
	LobbiesIndexed []string
	MaxCount       int
}

func NewManager(maxCount int) *LobbyManager {
	return &LobbyManager{
		Lobbies:  make(map[string]*lobby.Lobby),
		MaxCount: maxCount,
		Capacity: make(map[string]int),
	}
}

func (l *LobbyManager) ManagePlayer(player *models.Player, event *events.JoinEvent) (*lobby.Lobby, error) {

	if event.LobbyID == "" {
		// No Lobby was specified so the a lobby has to be searched for the player

		// Search for a existing lobby
		a := l.Search()

		log.Printf("Found lobbies %v", a)

		for _, id := range a {
			l.Lobbies[id].JoinPlayer(player)
			return l.Lobbies[id], nil
		}

		// No lobby fits so create a new one
		//
		id, err := l.CreateNew()
		if err == nil {
			l.Lobbies[id].JoinPlayer(player)
			return l.Lobbies[id], nil
		}

		log.Printf("Queueing in any full lobby")
		// get a random lobby int64 is faster than a intn
		i := rand.Int63() % int64(len(l.Lobbies))
		name := l.LobbiesIndexed[i]
		// Enqueue player in the lobby specific queue
		l.Lobbies[name].JoinPlayer(player)
		return l.Lobbies[name], nil
	}

	lobby, ok := l.Lobbies[event.ID]

	if !ok {
		return nil, errors.New("The specified LobbyId was invalid or the referenced lobby did not exist")
	}

	lobby.JoinPlayer(player)

	return lobby, nil
}
