package hand

import (
	"log"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/utils"
)

func (h *Hand) smallBlind() {

	//TODO better error handeling (3 tries then move to next)
	i := (h.Dealer + 1) % h.InCount
	utils.SendToAll(h.Players, events.NewWaitSmallBlindEvent(i))
	blind := h.In[i]
	e := <-blind.In

	log.Printf("event received %v", e)

	r, err := events.ToBlindSet(e)

	if err != nil {
		log.Printf("err %v", err)
	}

	h.Bank.SmallBlind(blind.ID, r)
	log.Printf("small blind received %v", r)
}

func (h *Hand) bigBlind() {
	//TODO better error handeling (3 tries then move to next)
	i := (h.Dealer + 2) % h.InCount
	utils.SendToAll(h.Players, events.NewWaitBigBlindEvent(i))
	blind := h.In[i]
	e := <-blind.In
	r, _ := events.ToBlindSet(e)
	h.Bank.BigBlind(blind.ID, r)
	log.Printf("big blind received %v", r)
}
