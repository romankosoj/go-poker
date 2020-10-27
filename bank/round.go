package bank

import (
	"errors"
	"log"
)

type Round struct {
	PlayerBets map[string]int
	Pot        int
	MaxBet     int
}

func NewRound(players map[string]int) *Round {
	bets := make(map[string]int)

	for n := range players {
		bets[n] = 0
	}

	return &Round{
		PlayerBets: bets,
		Pot:        0,
	}
}

func (r *Round) Conclude(winners []string) error {

	// transact pot to winners equaly

	share := r.Pot / len(winners)

	for _, n := range winners {
		log.Printf("User [%v] wins share %v ", n, share)
	}

	return nil
}

func (r *Round) PlayerBet(id string, amount int) error {
	bet, ok := r.PlayerBets[id]
	if !ok {
		return errors.New("Player not registered in bank (round)")
	}
	r.PlayerBets[id] = bet + amount
	r.Pot += amount
	if amount > r.MaxBet {
		r.MaxBet = amount
	}
	return nil
}
