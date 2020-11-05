package hand

import (
	"sync"

	"github.com/JohnnyS318/go-poker/bank"
	"github.com/JohnnyS318/go-poker/models"
	"github.com/JohnnyS318/go-poker/utils"
)

//Hand is one game of a session. it results in everybody but one folding or a showdown
type Hand struct {
	//Players includes all the Players who have started this hand. After a fold the player is still included
	Players         []models.Player
	Bank            *bank.Bank
	Board           [5]models.Card
	HoleCards       map[string][2]models.Card
	InCount         int
	Dealer          int
	Ended           bool
	cardGen         *utils.CardGenerator
	EndCallback     func(int)
	Blind           int
	bigBlindIndex   int
	smallBlindIndex int
	wg              sync.WaitGroup
}

//NewHand creates a new hand and sets the dealer to the next
func NewHand(bank *bank.Bank) *Hand {

	return &Hand{
		Bank:    bank,
		cardGen: utils.NewCardSelector(),
		Blind:   10,
	}
}

func (h *Hand) WhileNotEnded(f func()) {
	if !h.Ended {
		f()
	}
}
