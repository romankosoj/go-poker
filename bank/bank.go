package bank

import (
	"errors"
	"fmt"
	"log"

	"github.com/JohnnyS318/go-poker/models"
)

type Bank struct {
	PlayerValues map[string]int
	Round        *Round
}

func NewBank(players []models.Player) *Bank {
	values := make(map[string]int)

	for _, n := range players {
		values[n.ID] = n.BuyIn
	}

	return &Bank{
		PlayerValues: values,
		Round:        NewRound(values),
	}
}

func (b *Bank) ResetRound() {
	b.Round.Conclude()
	b.Round = NewRound(b.PlayerValues)
}

func (b *Bank) PlayerBet(id string, amount int) error {

	log.Printf("Player %v bets %d", id, amount)

	playerValue, ok := b.PlayerValues[id]

	if !ok {
		log.Printf("Player not registered in bank")
		return errors.New("Player not registered in bank")
	}

	if playerValue < amount {
		log.Printf("The player %v does not have the capacity to bet %v ", id, amount)
		return fmt.Errorf("The player does not have the capacity to bet %v ", amount)
	}

	if b.Round.MaxBet < 1 && amount < b.Round.MaxBet && playerValue != amount {

		// Player bet is les than round bet and is not an all in => invalid
		return errors.New("The player has to bet more or equal the round bet or do an all in")
	}

	log.Printf("Player can bet")

	//player can bet amount
	b.PlayerValues[id] = b.PlayerValues[id] - amount

	log.Printf("Subtracted amount")

	err := b.Round.PlayerBet(id, amount)

	log.Printf("Bet on the round")

	return err
}

func (b *Bank) SmallBlind(id string, amount int) error {
	return b.PlayerBet(id, amount)
}

func (b *Bank) BigBlind(id string, amount int) error {
	return b.PlayerBet(id, amount)
}

func (b *Bank) RemovePlayer(id string) error {
	_, ok := b.PlayerValues[id]
	if ok {
		delete(b.PlayerValues, id)
		return nil
	}
	return errors.New("Player not registered in bank")
}
