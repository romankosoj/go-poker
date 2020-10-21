package lobby

import (
	"errors"
	"log"
	"time"

	"github.com/JohnnyS318/go-poker/bank"
	"github.com/JohnnyS318/go-poker/hand"
	"github.com/JohnnyS318/go-poker/models"
)

type Lobby struct {
	Players      []models.Player `json:"players"`
	GameStarted  bool
	Callback     func()
	PlayerLeaves func(string) error
	MinBuyIn     int
	MaxBuyIn     int
}

func NewLobby() *Lobby {
	return &Lobby{
		Players: make([]models.Player, 0),
	}
}

func (l *Lobby) JoinPlayer(player *models.Player) error {
	log.Printf("Player joined lobby %v", player)

	if len(l.Players) <= 10 {
		l.Players = append(l.Players, *player)

		log.Printf("Player count in join lobby %v", len(l.Players))

		player.Out <- models.NewEvent("JOIN_LOBBY", "Yeah").ToRaw()

		return nil
	}

	return errors.New("The lobby is already full")
}

func (l *Lobby) RemovePlayerByID(id string) error {

	i := l.FindPlayerByID(id)

	if i < 0 {
		return errors.New("The player is not in the lobby")
	}

	log.Printf("Game runnnin ?: [%v]", l.GameStarted)

	if l.GameStarted {
		log.Printf("PlayerLeave exec")
		l.PlayerLeaves(id)
	}

	return l.RemovePlayer(i)
}

func (l *Lobby) RemovePlayer(i int) error {
	s := l.Players
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	l.Callback()
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

func (l *Lobby) Start() {
	l.GameStarted = true

	// SETUP
	dealer := -1
	for {
		time.Sleep(1 * time.Second)

		log.Printf("Game started")
		bank := bank.NewBank(l.Players)

		log.Printf("Bank created")

		hand := hand.NewHand(l.Players, bank, dealer)
		l.PlayerLeaves = func(id string) error {
			err := bank.RemovePlayer(id)
			if err != nil {
				return err
			}
			err = hand.PlayerLeaves(id)
			return err
		}
		dealer = hand.Start()
	}

}
