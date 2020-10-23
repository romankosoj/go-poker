package hand

import (
	"fmt"
	"log"
	"time"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/models"
	"github.com/JohnnyS318/go-poker/utils"
)

func (h *Hand) holeCards() {
	for i := range h.Players {
		if h.Players[i].Active {
			var cards [2]models.Card
			cards[0] = h.cardGen.SelectRandom()
			cards[1] = h.cardGen.SelectRandom()
			h.HoleCards[h.Players[i].ID] = cards
			log.Printf("Cards 0:%v 1:%v", cards[0].String(), cards[1].String())
			utils.SendToPlayer(&h.Players[i], events.NewHoleCardsEvent(cards))
		}
	}
}

func (h *Hand) recAction(blocking [10]*models.Player, i int, preflop bool) {

	if checkIfEmpty(blocking) {
		return
	}
	p := blocking[i]

	if p == nil || !p.Active {
		blocking[i] = nil
		h.recAction(blocking, (i+1)%10, preflop)
	}

	var payload int
	var a *events.Action
	success := false
	for j := 3; j > 0; j-- {
		a, err := h.waitForAction(i, preflop)
		if err != nil {
			h.playerError(i, fmt.Sprintf("The action was not valid. %v more tries", j))
		}

		if a.Action == events.FOLD {
			h.fold(p.ID)
			removeBlocking(blocking, i)
			success = true
			break
		}

		if !preflop && a.Action == events.CHECK {
			success = true

			_, i, err := h.searchByID(p.ID)
			if err == nil {
				addBlocking(blocking, i, p)
				break
			}
		}

		if a.Action == events.RAISE {
			if a.Payload > h.Bank.Round.MaxBet {
				err := h.Bank.PlayerBet(p.ID, a.Payload)
				if err == nil {
					success = true
					payload = a.Payload
					addAllButThisBlockgin(blocking, h.Players, p)

					break
				}
				h.playerError(i, fmt.Sprintf("Raise must be higher than the highest bet. %v more tries", j))
			}
		}

		if a.Action == events.BET {
			err := h.Bank.PlayerBet(p.ID, h.Bank.Round.MaxBet)
			if err == nil {
				success = true
				blocking[i] = nil
				payload = h.Bank.Round.MaxBet
				break
			}
			h.playerError(i, fmt.Sprintf("Bet must be equal to the current highest bet. %v more tries", j))
		}
	}

	if !success {
		a = &events.Action{
			Action:  events.FOLD,
			Payload: 0,
		}
		h.fold(p.ID)
		removeBlocking(blocking, i)
	}

	utils.SendToAll(h.Players, events.NewActionProcessedEvent(a.Action, payload, i, p))
	time.Sleep(3 * time.Second)

	if !checkIfEmpty(blocking) {
		h.recAction(blocking, (i+1)%10, preflop)
	}

	return

}

func (h *Hand) fold(id string) {
	i, err := h.searchByActiveID(id)

	log.Printf("folding player %v", id)

	if err != nil {
		log.Printf("Error during Folding Player[%v]: %v", i, err)
		return
	}

	h.Players[i].Active = false
	h.InCount--
}

func (h *Hand) playerError(i int, message string) {
	utils.SendToPlayerInList(h.Players, i, models.NewEvent("INVALID_ACTION", message))
}

func (h *Hand) actions(preflop bool) {
	var blocking [10]*models.Player
	for i, n := range h.Players {
		log.Printf("Blocking %v %v:", n.Active, n)
		if n.Active {
			blocking[i] = &n
		}
	}
	h.recAction(blocking, (h.Dealer+3)%h.InCount, preflop)
}

func (h *Hand) waitForAction(i int, preflop bool) (*events.Action, error) {

	if preflop {
		utils.SendToAll(h.Players, events.NewWaitForActionEvent(i, 0b111))
	} else {
		utils.SendToAll(h.Players, events.NewWaitForActionEvent(i, 0b1111))
	}

	n := h.Players[i]
	e := <-n.In
	action, err := events.ToAction(e)
	if err != nil {
		return nil, err
	}
	return action, nil
}

func (h *Hand) PlayerLeaves(id string) error {

	log.Printf("Player Leave invoked")

	n, i, err := utils.SearchByID(h.Players, id)

	if err != nil {
		return err
	}

	utils.SendToAll(h.Players, models.NewEvent(events.PLAYER_LEAVES, events.NewPlayerLeavesEvent(n, i)))

	if len(h.Players) < 3 {
		h.End()
	}

	return nil
}
