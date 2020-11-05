package hand

import (
	"log"
	"time"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/models"
	"github.com/JohnnyS318/go-poker/utils"
)

func (h *Hand) Start(players []models.Player, dealer int) {

	h.Dealer = dealer
	h.Players = players
	h.InCount = len(players)
	h.HoleCards = make(map[string][2]models.Card, len(players))

	//publish players and position

	time.Sleep(3 * time.Second)

	var publicPlayers []models.PublicPlayer

	for i := range h.Players {
		publicPlayers = append(publicPlayers, *h.Players[i].ToPublic())
	}

	for i := range h.Players {
		utils.SendToPlayerInList(h.Players, i, events.NewGameStartEvent(publicPlayers, i))
	}

	time.Sleep(3 * time.Second)

	// Publish choosen Dealer

	h.sendDealer()
	time.Sleep(3 * time.Second)

	//set predefined blinds
	h.setBlinds()

	log.Printf("Blinds Choosen")

	// Set players hole cards
	h.holeCards()

	log.Printf("Hole cards choosen")

	time.Sleep(3 * time.Second)

	h.actions(true)

	// Flop, turn and river are done here
	for i := 0; i < 5; i++ {
		h.Board[i] = h.cardGen.SelectRandom()
	}

	// send flop result

	utils.SendToAll(h.Players, events.NewFlopEvent(h.Board))

	//
	h.actions(false)
	// send turn result
	utils.SendToAll(h.Players, events.NewTurnEvent(h.Board))

	h.WhileNotEnded(func() {
		h.actions(false)
		// send river result
		utils.SendToAll(h.Players, events.NewRiverEvent(h.Board))
	})

	h.WhileNotEnded(func() {
		h.actions(false)
		log.Printf("Showdown")
	})

	winners := h.showdown()
	winningPlayers := make([]int, 0)
	for i := range winners {
		_, i, err := utils.SearchByID(h.Players, winners[i])
		if err == nil {
			winningPlayers = append(winningPlayers, i)
		}
	}

	winningPublic := make([]models.PublicPlayer, len(winningPlayers))
	for i, n := range winningPlayers {
		winningPublic[i] = publicPlayers[n]
	}

	share := h.Bank.ResetRound(winners)

	utils.SendToAll(h.Players, events.NewGameEndEvent(winningPublic, share))

	log.Printf("Hand Ended leaving now")

}

func (h *Hand) End() {
	log.Printf("Ending Hand due to error")
	h.Ended = true
}
