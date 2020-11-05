package lobby

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/JohnnyS318/go-poker/bank"
	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/hand"
	"github.com/JohnnyS318/go-poker/models"
	"github.com/JohnnyS318/go-poker/utils"
)

type Lobby struct {
	lock                 sync.RWMutex
	LobbyID              string          `json:"lobbyId"`
	Players              []models.Player `json:"players"`
	GameStarted          bool
	MinBuyIn             int
	MaxBuyIn             int
	Blinds               int
	ToBeRemoved          []int
	ToBeAdded            []*models.Player
	LobbyManagerCapacity map[string]int
	PlayerQueue          []*models.Player
	Bank                 *bank.Bank
	dealer               int
}

func NewLobby() *Lobby {
	return &Lobby{
		LobbyID:     GenerateLobbyID(),
		Players:     make([]models.Player, 0),
		ToBeRemoved: make([]int, 0),
		PlayerQueue: make([]*models.Player, 0),
		Bank:        bank.NewBank(),
		dealer:      -1,
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
	l.lock.Lock()
	l.GameStarted = true
	l.lock.Unlock()
	log.Printf("Lobby Started")
	// SETUP
<<<<<<< HEAD
	dealer := -1
=======
>>>>>>> e57f71a9581dbf0b0564fb148793e82dcc5f769a
	go func() {
		for len(l.Players) > 2 {
			log.Printf("Game started")
			time.Sleep(10 * time.Second)

			rand.Seed(time.Now().UnixNano())
			if l.dealer < 0 {
				l.dealer = rand.Intn(len(l.Players))
			} else {
				l.dealer = (l.dealer + 1) % len(l.Players)
			}
			log.Printf("Dealer choosen random %d", l.dealer)

			for i := range l.Players {
				l.Players[i].Active = true
			}

			hand := hand.NewHand(l.Players, l.Bank, l.dealer)

			hand.Start()

			l.lock.Lock()
			l.GameStarted = false
			l.lock.Unlock()
			if l.HasToBeRemoved() {
				l.RemoveAfterGame()
			}
			if l.HasToBeRemoved() {

				l.EmptyToBeAdded()
			}
		}
	}()
<<<<<<< HEAD
=======
}

func (l *Lobby) GetGameStarted() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.GameStarted
}

func (l *Lobby) HasToBeAdded() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.ToBeAdded) > 0
}

func (l *Lobby) HasToBeRemoved() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.ToBeRemoved) > 0
>>>>>>> e57f71a9581dbf0b0564fb148793e82dcc5f769a
}

func (l *Lobby) JoinPlayer(player *models.Player) {

	le := len(l.Players) + len(l.ToBeAdded)
	if le <= 10 {
		l.lock.Lock()
		l.ToBeAdded = append(l.ToBeAdded, player)
		l.lock.Unlock()

		if !l.GameStarted {
			l.EmptyToBeAdded()
		}

		utils.SendToPlayer(player, events.NewJoinSuccessEvent(l.LobbyID, l.Players, l.GameStarted, 0, le, l.MaxBuyIn, l.MinBuyIn, l.Blinds))

		l.Bank.AddPlayer(player)

		if len(l.Players) > 2 && !l.GameStarted {
			l.Start()
		}
	} else {
		l.EnqueuePlayer(player)
	}
}

func (l *Lobby) EmptyToBeAdded() {
	l.lock.Lock()
	defer l.lock.Unlock()
	for i := range l.ToBeAdded {
		if len(l.Players) < 10 {
			l.Players = append(l.Players, *l.ToBeAdded[i])
			log.Printf("Player joined lobby %v", l.ToBeAdded[i])
			log.Printf("Player count in join lobby %v", len(l.Players))
		} else {
			l.EnqueuePlayer(l.ToBeAdded[i])
		}
	}
	l.ToBeAdded = nil

}

func (l *Lobby) RemovePlayerByID(id string) error {

	i := l.FindPlayerByID(id)

	log.Printf("Removing player index %v", i)

	if i < 0 {
		return errors.New("The player is not in the lobby")
	}

	log.Printf("Game runing ?: [%v]", l.GameStarted)

	l.RemovePlayer(i)

	log.Printf("Players in lobby? %v", len(l.Players))

	return nil
}

func (l *Lobby) RemovePlayer(i int) {
	l.lock.Lock()
	l.ToBeRemoved = append(l.ToBeRemoved, i)
	gameStarted := l.GameStarted
	l.lock.Unlock()
	if !gameStarted {
		l.RemoveAfterGame()
	}
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
	l.lock.Lock()
	defer l.lock.Unlock()
	for _, i := range l.ToBeRemoved {
		l.Bank.RemovePlayer(l.Players[i].ID)
		l.Players = append(l.Players[:i], l.Players[i+1:]...)
		player, ok := l.DequeuePlayer()
		if ok {
			l.JoinPlayer(player)
		}
	}
}
