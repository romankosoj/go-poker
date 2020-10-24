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

func (h *Hand) recAction(blocking []int, i int, preflop bool) {

	blockingLength := len(blocking)
	// Check if blocking is an empty list
	if checkIfEmpty(blocking) {
		log.Printf("Anchor hit")
		return
	}

	k := blocking[i]

	if k < 0 || !h.Players[k].Active {
		log.Printf("player inactive")
		h.recAction(blocking, (i+1)%blockingLength, preflop)
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
			h.fold(h.Players[k].ID)
			blocking = removeBlocking(blocking, i)
			success = true
			succeededAction = events.Action{
				Action:  events.FOLD,w
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
			log.Printf("Raised to [%v] > [%v]", a.Payload, h.Bank.Round.MaxBet)
			if a.Payload > h.Bank.Round.MaxBet {
				amount := a.Payload
				err := h.Bank.PlayerBet(h.Players[k].ID, amount)
				if err == nil {
					success = true
					succeededAction = events.Action{
						Action:  events.RAISE,
						Payload: amount,
					}
					payload = amount
					addAllButThisBlockgin(blocking, h.Players, k)

					break
				}
				h.playerError(i, fmt.Sprintf("Raise must be higher than the highest bet. %v more tries", j))
			}
		}

		if a.Action == events.BET {
			err := h.Bank.PlayerBet(h.Players[k].ID, h.Bank.Round.MaxBet)
			if err == nil {
				success = true
				succeededAction = events.Action{
					Action:  a.Action,
					Payload: a.Payload,
				}
				blocking = removeBlocking(blocking, i)
				payload = h.Bank.Round.MaxBet
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
		h.fold(h.Players[k].ID)
		removeBlocking(blocking, i)
	}

	utils.SendToAll(h.Players, events.NewActionProcessedEvent(succeededAction.Action, payload, k))

	time.Sleep(1 * time.Second)

	if !checkIfEmpty(blocking) {
		log.Printf("Blocking is not empty")

		next := (i + 1) % len(blocking)
		if blockingLength > len(blocking) {
			next = i % len(blocking)
		}
		log.Printf("next is [%v] in %v", blocking[next], blocking)
		// blocking has changed now so the length is different and the
		h.recAction(blocking, next, preflop)
	}

	log.Printf("Done now leaving")

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
	log.Printf("Blocking: %v", blocking)
	log.Printf("Starting with [%v] bigBlind: [%v]", startIndexBlocking, h.bigBlindIndex)

	h.recAction(blocking, startIndexBlocking, preflop)
}

func (h *Hand) waitForAction(i int, preflop bool) (events.Action, error) {

	log.Printf("Waiting for Player Action from Player: %v", h.Players[i].String())

	if preflop {
		utils.SendToAll(h.Players, events.NewWaitForActionEvent(i, 0b111))
	} else {
		utils.SendToAll(h.Players, events.NewWaitForActionEvent(i, 0b1111))
	}
	e := <-h.Players[i].In
	action, err := events.ToAction(e)
	log.Printf("Action received %v with err %v", action, err)
	if err != nil {
		return events.Action{}, err
	}
	return *action, nil
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
