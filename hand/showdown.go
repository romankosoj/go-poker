package hand

import (
	"github.com/JohnnyS318/go-poker/models"
)

func (h *Hand) showdown() []string {

	if len(h.Players) < 1 {
		return nil
	}

	if len(h.Players) == 1 {
		//h.In[0] wins
		id := make([]string, 1)
		id[0] = h.Players[0].ID
		return id
	}

	ranks := make(map[string]int)
	for i := range h.Players {
		if h.Players[i].Active {
			cards := h.HoleCards[h.Players[i].ID]
			rank := EvaluatePlayer(append(cards[:], h.Board[:]...))
			ranks[h.Players[i].ID] = rank
		}
	}

	winners := make([]string, 0)

	highestRank := 0

	// Determin winner or winners
	for k, v := range ranks {
		if v == highestRank {
			// add to winner because rank is equal
			winners = append(winners, k)
		}
		if v > highestRank {
			//reset winners
			winners = nil
			highestRank = v
			winners = append(winners, k)
		}

	}

	return winners
}

//EvaluatePlayer generates a number as an identification of the players hole cards + the boards cards rank. it selects the best card seection and return the rank.
func EvaluatePlayer(cards []models.Card) int {

	maxRank := RankSpecificHand(cards[2:])
	for i := -1; i < 5; i++ {
		for j := -1; j < 5; j++ {
			if i == j {
				continue
			}

			//swap

			if i > -1 {
				cards[i+2], cards[0] = cards[0], cards[i+2]
			}

			if j > -1 {
				cards[j+2], cards[1] = cards[1], cards[j+2]
			}

			r := RankSpecificHand(cards[2:])

			if r > maxRank {
				maxRank = r
			}

			// swap back
			if i > -1 {
				cards[i+2], cards[0] = cards[0], cards[i+2]
			}

			if j > -1 {
				cards[j+2], cards[1] = cards[1], cards[j+2]
			}
		}
	}

	return maxRank
}

//RankSpecificHand generates a rank identificiation number for 5 card array out of the 7 cards.
func RankSpecificHand(cards []models.Card) int {
	//Hand identifier explanation:
	// 1 Byte Number:
	// 4 MSB: Describe Hand State (0: High Card - 8: Straight Flush)
	// 4 LSB: Describe State Identifier via Highest Card (0: Card Value 2 - 12: Card Value Ace)
	// Hand Identifier ranges from 0b00000000 (Theoretically Highest Card Two) to 0b10011101 (Straight Flush with Ace -> Royal Flush)
	// If a given Hand h1 is better than another h2, h1 > h2 always yields true.
	identifier := normalizeAce(cards[0].Value)

	sCol := normalizeAce(cards[0].Color)
	lP1 := -1
	lP2 := -1
	over := -1

	y := 0
	org := identifier

	m1 := -1
	m2 := -1

	min := identifier

	validStates := 0b111111111
	for i := 1; i < 5; i++ {
		open := normalizeAce(cards[i].Value)

		if normalizeAce(cards[i].Color) != sCol {
			validStates = validStates & 0b011011111
		}

		if org == open {
			y++
		} else {
			y += -1
		}

		if open == identifier || open == lP1 || open == lP2 || open == over {
			validStates = validStates & 0b011001110
			if m1 == -1 {
				m1 = open
			} else if open < m1 {
				m2 = open
			} else if open > m1 {
				m2 = m1
				m1 = open
			}
		} else {
			over = lP2
			if open > identifier {
				lP2 = lP1
				lP1 = identifier
				identifier = open
			} else if open > lP1 {
				lP2 = lP1
			} else if open > lP1 {
				lP2 = lP1
				lP1 = open
			} else if open > lP2 {
				lP2 = open
			} else {
				over = open
			}

			if over != -1 {
				validStates = validStates & 0b100100011
			}
		}

	}

	if lP2 != -1 {
		validStates = validStates & 0b101111111
	} else if y == 2 || y == -4 {
		validStates = validStates & 0b110111111
	}

	if m2 == -1 {
		validStates = validStates & 0b111111011
	} else {
		validStates = validStates & 0b111110111
	}

	if m1 == -1 && identifier-min == 4 {
		validStates = validStates & 0b100100000
	} else if m1 == -1 {
		validStates = validStates & 0b000100001
	}

	f := 0
	for f = 8; f >= 0; f-- {
		if (validStates & (1 << f)) != 0 {
			break
		}
	}

	if m1 == -1 {
		m1 = identifier
	}
	if m2 == -1 {
		if m1 == -1 {
			m2 = lP1
		} else {
			m2 = identifier
		}
	}

	return (f << 8) + (m1 << 4) + m2

}

func normalizeAce(number int) int {
	if number-1 < 0 {
		return 12
	}
	return number
}
