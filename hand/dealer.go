package hand

import (
	"errors"

	"github.com/JohnnyS318/go-poker/events"
	"github.com/JohnnyS318/go-poker/models"
	"github.com/JohnnyS318/go-poker/utils"
)

func (h *Hand) searchByID(id string) (*models.Player, int, error) {
	for i, n := range h.In {
		if n.ID == id {
			return &n, i, nil
		}
	}
	return nil, -1, errors.New("Player not in game")
}

func (h *Hand) sendDealer() {
	utils.SendToAll(h.In, models.NewEvent(events.DEALER_SET, h.Dealer))
}
