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
	lock                 sync.RWMutex
	LobbyID              string `json:"lobbyId"`
	Players              []models.Player
	PublicPlayers        []models.PublicPlayer `json:"players`
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
	hand                 *hand.Hand
}

func NewLobby() *Lobby {
	bank := bank.NewBank()
	return &Lobby{
		LobbyID:     GenerateLobbyID(),
		Players:     make([]models.Player, 0),
		ToBeRemoved: make([]int, 0),
		PlayerQueue: make([]*models.Player, 0),
		Bank:        bank,
		dealer:      -1,
		hand:        hand.NewHand(bank),
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

	log.Printf("Lobby Started")
	// SETUP
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

			l.lock.Lock()
			l.GameStarted = true
			l.lock.Unlock()
			l.hand.Start(l.Players, l.dealer)
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
}

func (l *Lobby) JoinPlayer(player *models.Player) {

	le := len(l.Players) + len(l.ToBeAdded)
	if le <= 10 {
		l.lock.Lock()
		l.ToBeAdded = append(l.ToBeAdded, player)
		l.lock.Unlock()
		gameStarted := l.GetGameStarted()
		if !gameStarted {
			l.EmptyToBeAdded()
		}
		l.Bank.AddPlayer(player)

		time.Sleep(200 * time.Millisecond)
		log.Printf("Join Success event to [%v]", player.String())
		utils.SendToPlayer(player, events.NewJoinSuccessEvent(l.LobbyID, l.PublicPlayers, gameStarted, 0, le, l.MaxBuyIn, l.MinBuyIn, l.Blinds))

		if len(l.Players) > 2 && !gameStarted {
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
			for j := range l.Players {
				if l.Players[j].ID == l.ToBeAdded[i].ID {
					continue
				}
			}
			utils.SendToAll(l.Players, events.NewPlayerJoinEvent(l.ToBeAdded[i].ToPublic(), len(l.Players)-1))
			l.Players = append(l.Players, *l.ToBeAdded[i])
			l.PublicPlayers = append(l.PublicPlayers, *l.ToBeAdded[i].ToPublic())
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

	if i < 0 {
		return errors.New("The player is not in the lobby")
	}

	l.lock.Lock()
	l.ToBeRemoved = append(l.ToBeRemoved, i)
	l.lock.Unlock()
	l.hand.Fold(id)
	if !l.GetGameStarted() {
		l.RemoveAfterGame()
	}

	log.Printf("Players in lobby? %v", len(l.Players))

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
	l.lock.Lock()
	defer l.lock.Unlock()
	for _, i := range l.ToBeRemoved {
		log.Printf("Removing player [%v] in list with length %v", i, len(l.Players))
		if len(l.Players) > i {
			l.Bank.RemovePlayer(l.Players[i].ID)
			l.Players = append(l.Players[:i], l.Players[i+1:]...)
			l.PublicPlayers = append(l.PublicPlayers[:i], l.PublicPlayers[i+1:]...)
			player, ok := l.DequeuePlayer()
			if ok {
				l.JoinPlayer(player)
			}
		}
	}
}
