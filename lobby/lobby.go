package lobby

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/JohnnyS318/go-poker/bank"
	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/hand"
	"github.com/JohnnyS318/go-poker/models"
	"github.com/JohnnyS318/go-poker/utils"
)

type Lobby struct {
	LobbyID              string          `json:"lobbyId"`
	Players              []models.Player `json:"players"`
	GameStarted          bool
	MinBuyIn             int
	MaxBuyIn             int
	Blinds               int
	ToRemove             []int
	LobbyManagerCapacity map[string]int
	PlayerQueue          []*models.Player
	Bank                 *bank.Bank
}

func NewLobby() *Lobby {
	return &Lobby{
		LobbyID:     GenerateLobbyID(),
		Players:     make([]models.Player, 0),
		ToRemove:    make([]int, 0),
		PlayerQueue: make([]*models.Player, 0),
		Bank:        bank.NewBank(),
	}
}

func (l *Lobby) EnqueuePlayer(player *models.Player) {
	l.PlayerQueue = append(l.PlayerQueue, player)
}

func (l *Lobby) DequeuePlayer() (player *models.Player, ok bool) {
	if len(l.PlayerQueue) > 0 {
		if len(l.Players) < 10 {
			player := l.PlayerQueue[0]
			l.Players = append(l.Players, *player)
			l.PlayerQueue = l.PlayerQueue[1:]
			return player, true
		}
	}

	return nil, false
}

func (l *Lobby) Start() {
	l.GameStarted = true
	log.Printf("Lobby Started")
	// SETUP
	dealer := -1

	for len(l.Players) > 2 {
		log.Printf("Game started")
		time.Sleep(10 * time.Second)

		var wg sync.WaitGroup

		rand.Seed(time.Now().UnixNano())
		if dealer < 0 {
			dealer = rand.Intn(len(l.Players))
		} else {
			dealer = (dealer + 1) % len(l.Players)
		}
		log.Printf("Dealer choosen random %d", dealer)

		for i := range l.Players {
			l.Players[i].Active = true
		}

		hand := hand.NewHand(l.Players, l.Bank, dealer)

		hand.Start()

		l.GameStarted = false
		if len(l.ToRemove) < 1 {
			l.RemoveAfterGame()
		}

	}
}

func (l *Lobby) JoinPlayer(player *models.Player) {
	if len(l.Players) <= 10 {
		i := len(l.Players)
		l.Players = append(l.Players, *player)

		log.Printf("Player joined lobby %v", player)
		log.Printf("Player count in join lobby %v", len(l.Players))

		utils.SendToPlayer(player, events.NewJoinSuccessEvent(l.LobbyID, l.Players, l.GameStarted, 0, i, l.MaxBuyIn, l.MinBuyIn, l.Blinds))

		l.Bank.AddPlayer(player)

		if len(l.Players) > 2 && !l.GameStarted {
			l.Start()
		}
	} else {
		l.EnqueuePlayer(player)
	}
}

func (l *Lobby) RemovePlayerByID(id string) error {

	i := l.FindPlayerByID(id)

	log.Printf("Removing player index %v", i)

	if i < 0 {
		return errors.New("The player is not in the lobby")
	}

	log.Printf("Game runing ?: [%v]", l.GameStarted)

	err := l.RemovePlayer(i)

	log.Printf("Players in lobby? %v", len(l.Players))

	return err
}

func (l *Lobby) RemovePlayer(i int) error {
	if l.GameStarted {
		l.ToRemove = append(l.ToRemove, i)
	} else {
		//Keep order so that the next dealer is choosen correctly.
		l.Players = append(l.Players[:i], l.Players[i+1:]...)

		player, ok := l.DequeuePlayer()
		if ok {
			l.JoinPlayer(player)
		}

		// Non blocking channel send (possibly nobody was listening for a player leaves and we have to continue)
		//select {
		//case l.PlayerLeavesChannel <- l.LobbyID:
		//}
	}
	return nil
}

func (l *Lobby) FindPlayerByID(id string) int {
	for i, n := range l.Players {
		if n.ID == id {
			return i
		}
	}
	return -1
}

func (l *Lobby) HasCapacaty() bool {
	return len(l.Players) < 10
}

//RemoveAfterGame removes the left players from the lobby after a game has finished. During a game the player is counted as folded.
func (l *Lobby) RemoveAfterGame() {
	for _, i := range l.ToRemove {
		l.Bank.RemovePlayer(l.Players[i].ID)
		l.Players = append(l.Players[:i], l.Players[i+1:]...)
		player, ok := l.DequeuePlayer()
		if ok {
			l.JoinPlayer(player)
		}
	}
}
