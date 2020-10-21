package hand

import (
	"log"
	"time"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/models"
	"github.com/JohnnyS318/go-poker/utils"
)

func (h *Hand) Start() int {

	//publish players and position

	var publicPlayer []models.PublicPlayer

	for _, n := range h.In {
		publicPlayer = append(publicPlayer, *n.ToPublic())
	}

	for i := range h.In {
		utils.SendToPlayerInList(h.In, i, events.NewGameStartEvent(publicPlayer, i))
	}

	time.Sleep(5 * time.Second)

	// Publish choosen Dealer

	h.sendDealer()

	time.Sleep(5 * time.Second)

	// choose blinds
	h.smallBlind()
	h.bigBlind()

	log.Printf("Blinds Choosen")

	// Set players hole cards
	h.holeCards()

	log.Printf("Hole cards choosen")

	time.Sleep(5 * time.Second)

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

	h.actions(false)

	// send river result

	utils.SendToAll(h.Players, events.NewRiverEvent(h.Board))

	h.actions(false)

	log.Printf("Showdown")

	h.showdown()

	return h.Dealer
}
