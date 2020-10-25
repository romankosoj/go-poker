package hand

import (
	"github.com/JohnnyS318/go-poker/models"
)

func (h *Hand) showdown() {

	if len(h.Players) < 1 {
		return
	}

	if len(h.Players) == 1 {
		//h.In[0] wins
	}
}

func (h *Hand) RankSpecificHand(cards []models.Card) uint16{
	//Hand identifier explanation:
	// 1 Byte Number:
	// 4 MSB: Describe Hand State (0: High Card - 8: Straight Flush)
	// 4 LSB: Describe State Identifier via Highest Card (0: Card Value 2 - 12: Card Value Ace)
	// Hand Identifier ranges from 0b00000000 (Theoretically Highest Card Two) to 0b10011101 (Straight Flush with Ace -> Royal Flush)
	// If a given Hand h1 is better than another h2, h1 > h2 always yields true.
	identifier := normalizeAce(cards[0].Value)

	sCol := normalizeAce(cards[0].Color);
    lP1 := -1;
    lP2 := -1;
	over  := -1
	
    //Possible 9 hand states are one by one excluded by masking a 9 bit number.
	validStates := 0b111111111

	for i := 0; i < 5; i++ {
		open := normalizeAce(cards[i].Value)
		r := open - identifier;
		if r > 4 || r < 4{
			validStates = validStates & 0b011101111
		}
		if open == identifier || open == lP1 || open == lp2 || open == over {
			validStates = validStates &  0b011001110;
		} else {
			over = lP2
			if open > identifier{
                lP2 = lP1
                lP1 = identifier
                identifier = open
			} else if open > lp1{
				lP2 = lP1
				lP1 = open
			} else if open > lp2 {
				
			}
		}

	}

}

func normalizeAce(number int) int {
	if number-1 < 0 {
		return 12
	}
	return number
}

func v(n int, c []models.Card) int {
	return normalizeAce(c[n].Value)
})