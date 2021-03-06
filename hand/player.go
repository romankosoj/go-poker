package hand

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/models"
	"github.com/JohnnyS318/go-poker/utils"
)

func (h *Hand) actions(preflop bool) {

	var startIndexPlayers int
	for j := 1; j <= len(h.Players); j++ {
		startIndexPlayers = (h.bigBlindIndex + j) % len(h.Players)
		if h.Players[startIndexPlayers].Active {
			break
		}
	}

	startIndexBlocking := -1
	blocking := make([]int, 0)
	for i, n := range h.Players {
		if n.Active {
			blocking = append(blocking, i)
			if startIndexPlayers == i {
				startIndexBlocking = i
			}
		}
	}

	h.recAction(blocking, startIndexBlocking%len(blocking), preflop)
}

func (h *Hand) recAction(blocking []int, i int, preflop bool) {

	if h.InCount < 2 {
		return
	}

	removed := false
	// Check if blocking is an empty list
	if checkIfEmpty(blocking) {
		log.Printf("Anchor hit")
		return
	}

	k := blocking[i]

	if k < 0 || !h.Players[k].Active {
		log.Printf("player inactive")

		// remove from blocking list

		h.recAction(blocking, (i+1)%len(blocking), preflop)
		return
	}

	payload := 0
	var succeededAction events.Action
	success := false
	for j := 3; j > 0; j-- {
		a, err := h.waitForAction(k, preflop)
		succeededAction = a
		if err != nil {
			h.playerError(i, fmt.Sprintf("The action was not valid. %v more tries", j))
		}

		if a.Action == events.FOLD {
			h.Fold(h.Players[k].ID)
			blocking = removeBlocking(blocking, i)
			removed = true
			success = true
			succeededAction = events.Action{
				Action:  events.FOLD,
				Payload: a.Payload,
			}
			break
		}

		if !preflop && a.Action == events.CHECK {
			success = true
			succeededAction = events.Action{
				Action:  events.CHECK,
				Payload: a.Payload,
			}
			if err == nil {
				addBlocking(blocking, k)
				break
			}
		}

		if a.Action == events.RAISE {
			max := h.Bank.GetMaxBet()
			log.Printf("Raised to [%v] > [%v]", a.Payload, max)
			if a.Payload > max {
				amount := a.Payload
				err := h.Bank.PlayerBet(h.Players[k].ID, amount)
				if err == nil {
					success = true
					succeededAction = events.Action{
						Action:  events.RAISE,
						Payload: amount,
					}
					payload = amount
					blocking = addAllButThisBlockgin(blocking, h.Players, k)
					removed = true
					break
				}
				h.playerError(i, fmt.Sprintf("Raise must be higher than the highest bet. %v more tries", j))
			}
		}

		if a.Action == events.BET {
			max := h.Bank.GetMaxBet()
			err := h.Bank.PlayerBet(h.Players[k].ID, max)
			if err == nil {
				success = true
				succeededAction = events.Action{
					Action:  a.Action,
					Payload: a.Payload,
				}
				blocking = removeBlocking(blocking, i)
				removed = true
				payload = max
				break
			}
			h.playerError(i, fmt.Sprintf("Bet must be equal to the current highest bet. %v more tries", j))
		}
	}

	if !success {
		log.Printf("Failed Processing action. sending now")
		succeededAction = events.Action{
			Action:  events.FOLD,
			Payload: 0,
		}
		h.Fold(h.Players[k].ID)
		removeBlocking(blocking, i)
		removed = true
	}

	utils.SendToAll(h.Players, events.NewActionProcessedEvent(succeededAction.Action, payload, k, h.Bank.GetTotalPlayerBet(h.Players[k].ID)))

	time.Sleep(1 * time.Second)

	if !checkIfEmpty(blocking) {

		next := (i + 1) % len(blocking)
		if removed {
			next = i % len(blocking)
		}
		// blocking has changed now so the length is different and the
		h.recAction(blocking, next, preflop)
	}

	return

}

func (h *Hand) holeCards() {
	for i := range h.Players {
		if h.Players[i].Active {
			var cards [2]models.Card
			cards[0] = h.cardGen.SelectRandom()
			cards[1] = h.cardGen.SelectRandom()
			h.HoleCards[h.Players[i].ID] = cards
			utils.SendToPlayer(&h.Players[i], events.NewHoleCardsEvent(cards))
		}
	}
}

func (h *Hand) Fold(id string) error {
	i, err := h.searchByActiveID(id)

	if err != nil {
		return err
	}

	if i < 0 || i >= len(h.Players) {
		return errors.New("Something went wrong")
	}
	h.Players[i].Active = false
	h.InCount--
	utils.SendToAll(h.Players, events.NewActionProcessedEvent(events.FOLD, 0, i, h.Bank.GetTotalPlayerBet(h.Players[i].ID)))
	log.Printf("Folded player [%v]", h.Players[i].String())
	return nil
}

func (h *Hand) playerError(i int, message string) {
	utils.SendToPlayerInList(h.Players, i, models.NewEvent("INVALID_ACTION", message))
}

func (h *Hand) waitForAction(i int, preflop bool) (events.Action, error) {
	if preflop {
		utils.SendToAll(h.Players, events.NewWaitForActionEvent(i, 0b111))
	} else {
		utils.SendToAll(h.Players, events.NewWaitForActionEvent(i, 0b1111))
	}
	e := <-h.Players[i].In
	action, err := events.ToAction(e)
	if err != nil {
		return events.Action{}, err
	}
	return *action, nil
}

func (h *Hand) PlayerLeaves(id string) error {

	log.Printf("Player Leave invoked")

	_, i, err := utils.SearchByID(h.Players, id)

	if err != nil {
		return err
	}

	err = h.Fold(id)

	if err != nil {
		return err
	}

	utils.SendToAll(h.Players, events.NewActionProcessedEvent(events.FOLD, 0, i, h.Bank.GetTotalPlayerBet(h.Players[i].ID)))

	if len(h.Players) < 2 {
		h.End()
	}

	return nil
}
