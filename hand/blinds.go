package hand

import (
	"errors"
	"log"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/utils"
)

func (h *Hand) setBlinds() error {
	success := false
	s := 0

	var smallBlindIndex int
	for !success {

		if s >= len(h.Players)-1 {
			return errors.New("Nobody can bet smallBlind and bigBlind")
		}

		smallBlindAmount := h.Blind / 2

		log.Printf("Deciding small blind")

		var i int
		for j := 1; j < len(h.Players)-1; j++ {
			i = (h.Dealer + j) % len(h.Players)
			if h.Players[i].Active && h.Dealer != i {
				break
			}
		}

		log.Printf("Player [%v] has to bet [%v]: %v ", i, smallBlindAmount, h.Players[i].String())

		err := h.Bank.PlayerBet(h.Players[i].ID, smallBlindAmount)

		if err != nil {
			log.Printf("Folding player due to invalid buyin in smallBlind %v", err)
			h.Fold(h.Players[i].ID)
			s++
			continue
		}

		log.Printf("Bet succeeded sending now.")
		utils.SendToAll(h.Players, events.NewActionProcessedEvent(2, smallBlindAmount, i, h.Bank.GetTotalPlayerBet(h.Players[i].ID)))
		success = true
		smallBlindIndex = i
		s = 0
		log.Printf("sucess [%v]", smallBlindIndex)
	}

	success = false
	for !success {

		if s >= len(h.Players) {
			return errors.New("Nobody can bet bigBlind")
		}

		var i int
		for j := 1; j <= len(h.Players); j++ {
			i = (smallBlindIndex + j) % len(h.Players)
			if h.Players[i].Active && h.Dealer != i && smallBlindIndex != i {
				break
			}
		}

		log.Printf("Player [%v] has to bet [%v]: %v ", i, h.Blind, h.Players[i].String())

		err := h.Bank.PlayerBet(h.Players[i].ID, h.Blind)

		if err != nil {
			log.Printf("Folding player due to invalid buyin in bigBlind")
			h.Fold(h.Players[i].ID)
			s++
			continue
		}
		log.Printf("Bet succeeded sending now.")
		utils.SendToAll(h.Players, events.NewActionProcessedEvent(2, h.Blind, i, h.Bank.GetTotalPlayerBet(h.Players[i].ID)))
		success = true
		h.bigBlindIndex = i
		log.Printf("success [%v]", h.bigBlindIndex)
	}

	return nil
}
