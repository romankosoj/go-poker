package lobbies

import (
	"log"
	"time"

	"github.com/JohnnyS318/go-poker/lobby"
	"github.com/JohnnyS318/go-poker/models"
)

type LobbyManager struct {
	Lobbies  []lobby.Lobby `json:"lobbies"`
	MaxCount int
}

func NewManager(maxCount int) *LobbyManager {
	return &LobbyManager{
		Lobbies:  make([]lobby.Lobby, 0),
		MaxCount: maxCount,
	}
}

func (l *LobbyManager) ManagePlayer(player *models.Player) (lobby.Lobby, int, error) {

	// Search for a existing lobby
	a := l.Search()

	log.Printf("found lobbies %v", a)

	for _, i := range a {
		lob, ind, err := l.JoinPlayer(i, player)
		if err == nil {
			return lob, ind, nil
		}
		log.Printf("Error during joining %v", err)
	}

	c := make(chan int)

	// No lobby fits so create a new one
	//
	i := l.CreateNew(c)

	if i >= 0 {
		return l.JoinPlayer(i, player)
	}

	log.Printf("Waiting for empty lobby")

	// No new lobby can be crated, due to this we have to wait until a player leaves a given lobby.
	// We register in a single line queue to
	i = <-c
	return l.JoinPlayer(i, player)
}

func (l *LobbyManager) JoinPlayer(i int, player *models.Player) (lobby.Lobby, int, error) {

	log.Printf("Joining lobby")

	err := l.Lobbies[i].JoinPlayer(player)
	if err != nil {
		return lobby.Lobby{}, -1, err
	}

	log.Printf("Player count %v", len(l.Lobbies[i].Players))

	if len(l.Lobbies[i].Players) > 2 && !l.Lobbies[i].GameStarted {

		time.Sleep(1 * time.Second)

		log.Printf("Starting game now")

		// // create callback that game finished
		// g := game.Create(l.Lobbies[i].Players, func() {
		// 	l.Lobbies[i].GameStarted = false
		// })

		// // create notification that a player left the lobby
		// l.Lobbies[i].PlayerLeaves = g.PlayerLeaves

		// // Start the game on different go routine (thread) so that multple games can be run simultaneously.
		// go g.Start()

		go l.Lobbies[i].Start()
	}
	return l.Lobbies[i], i, nil
}
