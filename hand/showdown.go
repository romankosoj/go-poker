package hand

import "github.com/JohnnyS318/go-poker/models"

func (h *Hand) showdown() {

	if len(h.Players) == 1 {
		//h.In[0] wins
	}

	//
	counts := make(map[string]int, len(h.Players))

	for i, n := range h.Players {
		hole := h.HoleCards[n.ID]
		cards := append(hole[:], h.Board[:]...)
		counts[n.ID] = h.RankPlayer(i, &n, cards[:7])
	}

}

func (h *Hand) RankPlayer(i int, p *models.Player, cards []models.Card) int {

	highestCard := -1
	dColor := make(map[int]int, 4)

	// k: card value
	// v: count
	dValue := make(map[int]int, 13)

	for _, n := range cards {
		count, ok := dValue[n.Value]
		if ok {
			dValue[n.Value] = count + 1
		} else {
			dValue[n.Value] = 1
		}

		count, ok = dColor[n.Color]
		if ok {
			dColor[n.Color] = count + 1
		} else {
			dColor[n.Color] = 1
		}

		if n.Value > highestCard {
			highestCard = n.Value
		}
	}

	// [1:2, 2:2, 3:0, 4:1, ...]
	return 1
}
