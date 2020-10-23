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

		var i int
		for j := 1; j < len(h.Players)-1; j++ {
			i = (h.Dealer + j) % len(h.Players)
			if h.Players[i].Active && h.Dealer != i {
				break
			}
		}

		player := &h.Players[i]

		err := h.Bank.PlayerBet(player.ID, smallBlindAmount)

		if err != nil {
			log.Printf("Folding player due to invalid buyin in smallBlind")
			h.fold(player.ID)
			continue
		}

		utils.SendToAll(h.Players, events.NewActionProcessedEvent(2, smallBlindAmount, i))
		success = true
		smallBlindIndex = i
	}

	success = false
	for !success {
		var i int
		for j := 1; j < len(h.Players)-2; j++ {
			i = (smallBlindIndex + j) % len(h.Players)
			if h.Players[i].Active && h.Dealer != i && smallBlindIndex != i {
				break
			}
		}

		player := &h.Players[i]
		err := h.Bank.PlayerBet(player.ID, h.Blind)

		if err != nil {
			log.Printf("Folding player due to invalid buyin in bigBlind")
			h.fold(player.ID)
			continue
		}
		utils.SendToAll(h.Players, events.NewActionProcessedEvent(2, h.Blind, i))
		success = true
	}
}
