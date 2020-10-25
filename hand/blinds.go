package hand

import (
	"log"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/utils"
)

func (h *Hand) setBlinds() {
	success := false

	var smallBlindIndex int
	for !success {
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
			log.Printf("Folding player due to invalid buyin in smallBlind")
			h.fold(h.Players[i].ID)
			continue
		}

		log.Printf("Bet succeeded sending now.")
		utils.SendToAll(h.Players, events.NewActionProcessedEvent(2, smallBlindAmount, i))
		success = true
		smallBlindIndex = i
		log.Printf("sucess [%v]", smallBlindIndex)
	}

	success = false
	for !success {
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
			h.fold(h.Players[i].ID)
			continue
		}
		log.Printf("Bet succeeded sending now.")
		utils.SendToAll(h.Players, events.NewActionProcessedEvent(2, h.Blind, i))
		success = true
		h.bigBlindIndex = i
		log.Printf("success [%v]", h.bigBlindIndex)
	}
}
