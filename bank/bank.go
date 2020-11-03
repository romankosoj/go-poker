package bank

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/JohnnyS318/go-poker/models"
)

type Bank struct {
	lock         sync.RWMutex
	PlayerWallet map[string]int

	PlayerBets map[string]int
	Pot        int
	MaxBet     int
}

func NewBank() *Bank {
	return &Bank{
		PlayerWallet: make(map[string]int),
		PlayerBets:   make(map[string]int),
	}
}

func (b *Bank) AddPlayer(player *models.Player) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.PlayerWallet[player.ID] = player.BuyIn
	b.PlayerBets[player.ID] = 0
}

func (b *Bank) ResetRound(winners []string) int {
	b.lock.Lock()
	defer b.lock.Unlock()

	if winners == nil || b.Pot == 0 {
		return -1
	}

	share := b.Pot / len(winners)

	for _, n := range winners {
		b.PlayerWallet[n] += share
		log.Printf("User [%v] wins share %d", n, share)
	}

	b.Pot = 0
	return share
}

func (b *Bank) PlayerBet(id string, amount int) error {

	b.lock.Lock()
	defer b.lock.Unlock()

	log.Printf("Player %v bets %d", id, amount)
	playerValue, ok := b.PlayerWallet[id]

	if !ok {
		log.Printf("Player not registered in bank")
		return errors.New("Player not registered in bank")
	}

	if playerValue < amount {
		log.Printf("The player %v does not have the capacity to bet %v ", id, amount)
		return fmt.Errorf("The player does not have the capacity to bet %v ", amount)
	}

	if amount < b.MaxBet && playerValue != amount {
		// Player bet is les than round bet and is not an all in => invalid
		return errors.New("The player has to bet more or equal the round bet or do an all in")
	}

	//player can bet amount
	b.PlayerWallet[id] = b.PlayerWallet[id] - amount
	b.PlayerBets[id] = amount

	log.Printf("Bet success")

	return nil
}

func (b *Bank) RemovePlayer(id string) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	_, ok := b.PlayerWallet[id]
	if ok {
		delete(b.PlayerWallet, id)
		return nil
	}
	return errors.New("Player not registered in bank")
}

func (b *Bank) GetPot() int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.Pot
}

func (b *Bank) GetMaxBet() int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.MaxBet
}

func (b *Bank) GetPlayerWallet() map[string]int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.PlayerWallet
}

func (b *Bank) GetPlayerBets() map[string]int {
	b.lock.RLock()
	defer b.lock.RUnlock()
	return b.PlayerBets
}
