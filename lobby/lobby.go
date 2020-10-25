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
	ToRemove     []int
}

func NewLobby() *Lobby {
	return &Lobby{
		Players:  make([]models.Player, 0),
		ToRemove: make([]int, 0),
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

	log.Printf("Removing player index %v", i)

	if i < 0 {
		return errors.New("The player is not in the lobby")
	}

	log.Printf("Game runing ?: [%v]", l.GameStarted)

	if l.GameStarted {
		log.Printf("PlayerLeave exec")
		l.PlayerLeaves(id)
	}

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
		l.Callback()
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

func (l *Lobby) Start() {
	l.GameStarted = true

	// SETUP
	dealer := -1
	for len(l.Players) > 2 {
		log.Printf("Game started")

		for i := range l.Players {
			l.Players[i].Active = true
		}

		time.Sleep(1 * time.Second)

		bank := bank.NewBank(l.Players)

		hand := hand.NewHand(l.Players, bank, dealer)
		l.PlayerLeaves = func(id string) error {
			err := hand.PlayerLeaves(id)
			return err
		}
		dealer = hand.Start()
		l.GameStarted = false
		if len(l.ToRemove) < 1 {
			l.RemoveAfterGame(bank)
		}

	}
}

//RemoveAfterGame removes the left players from the lobby after a game has finished. During a game the player is counted as folded.
func (l *Lobby) RemoveAfterGame(bank *bank.Bank) {
	for _, i := range l.ToRemove {
		bank.RemovePlayer(l.Players[i].ID)
		l.Players = append(l.Players[:i], l.Players[i+1:]...)
		l.Callback()
	}
}
