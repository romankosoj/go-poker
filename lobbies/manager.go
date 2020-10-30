package lobbies

import (
	"errors"
	"log"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/lobby"
	"github.com/JohnnyS318/go-poker/models"
)

type LobbyManager struct {
	Lobbies     map[string]*lobby.Lobby
	Capacity    map[string]int
	MaxCount    int
	PlayerQueue *PlayerQueue
}

func NewManager(maxCount int) *LobbyManager {
	return &LobbyManager{
		Lobbies:     make(map[string]*lobby.Lobby),
		MaxCount:    maxCount,
		Capacity:    make(map[string]int),
		PlayerQueue: NewPlayerQueue(),
	}
}

func (l *LobbyManager) ManagePlayer(player *models.Player, event *events.JoinEvent) (*lobby.Lobby, error) {

	if event.LobbyID == "" {
		// No Lobby was specified so the a lobby has to be searched for the player

		// Search for a existing lobby
		a := l.Search()

		log.Printf("Found lobbies %v", a)

		for _, id := range a {
			lobby := l.Lobbies[id]
			err := l.Join(lobby, player)
			if err == nil {
				return lobby, nil
			}
			log.Printf("Error during lobby joining %v", err)
		}

		c := make(chan string)

		// No lobby fits so create a new one
		//
		id, err := l.CreateNew(c)

		if err == nil {
			lobby := l.Lobbies[id]
			err := l.Join(lobby, player)

			if err == nil {
				return lobby, nil
			}
		}

		log.Printf("Waiting for empty lobby")

		// No new lobby can be crated, due to this we have to wait until a player leaves a given lobby.
		// We register in a single line queue to
		l.PlayerQueue.Enqueue(player, event.ID)
		id = <-c

		lobby := l.Lobbies[id]

		err = l.Join(lobby, player)

		if err != nil {
			return nil, errors.New("Joining after waiting was unsuccessfull")
		}
		return lobby, nil

	}
	lobby, ok := l.Lobbies[event.ID]

	if !ok {
		return nil, errors.New("The specified LobbyId was invalid or the referenced lobby did not exist")
	}

	if lobby.HasCapacaty() {
		err := l.Join(lobby, player)

		if err == nil {
			return lobby, nil
		}
	}

	return nil, errors.New("Something went wrong")
}

func (l *LobbyManager) Join(lobby *lobby.Lobby, player *models.Player) error {
	log.Printf("Joining Lobby [%v]", lobby.LobbyID)
	err := lobby.JoinPlayer(player)

	if err != nil {
		return err
	}

	log.Printf("Player count %v", len(lobby.Players))
	return nil
}
