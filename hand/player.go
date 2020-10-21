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
	for _, n := range h.In {
		var cards [2]models.Card
		cards[0] = h.cardGen.SelectRandom()
		cards[1] = h.cardGen.SelectRandom()
		h.HoleCards[n.ID] = cards
		log.Printf("Cards 0:%v 1:%v", cards[0].String(), cards[1].String())
		utils.SendToPlayer(&n, events.NewHoleCardsEvent(cards))
	}
}

func (h *Hand) recAction(blocking [10]*models.Player, i int, preflop bool) {

	p := blocking[i]

	if p == nil {
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
					addAllButThisBlockgin(blocking, h.In, p)

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

	utils.SendToAll(h.In, events.NewActionProcessedEvent(a.Action, payload, i, p))
	time.Sleep(3 * time.Second)

	if !checkIfEmpty(blocking) {
		h.recAction(blocking, (i+1)%10, preflop)
	}

	return

}

func (h *Hand) fold(id string) {
	_, i, err := h.searchByID(id)
	if err != nil {
		log.Printf("Error during Folding Player[%v]: %v", i, err)
		return
	}
	h.In = append(h.In[:i], h.In[i+1:]...)
	h.InCount--
}

func (h *Hand) playerError(i int, message string) {
	utils.SendToPlayerInList(h.In, i, models.NewEvent("INVALID_ACTION", message))
}

func (h *Hand) actions(preflop bool) {
	var blocking [10]*models.Player
	for i, n := range h.In {
		blocking[i] = &n
	}
	h.recAction(blocking, (h.Dealer+3)%h.InCount, preflop)
}

func (h *Hand) waitForAction(i int, preflop bool) (*events.Action, error) {

	utils.SendToAll(h.In, events.NewWaitForActionEvent(i, preflop))

	n := h.In[i]
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

	h.Players = append(h.Players[:i], h.Players[i+1:]...)
	h.In = append(h.In[:i], h.In[i+1:]...)

	utils.SendToAll(h.Players, models.NewEvent(events.PLAYER_LEAVES, events.NewPlayerLeavesEvent(n, i)))

	return nil
}
