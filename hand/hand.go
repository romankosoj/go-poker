package hand

import (
	"log"
	"math/rand"
	"time"

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
}

//NewHand creates a new hand and sets the dealer to the next
func NewHand(players []models.Player, bank *bank.Bank, previousDealer int) *Hand {

	rand.Seed(time.Now().UnixNano())
	var dealer int
	if previousDealer < 0 {
		l := rand.Intn(len(players))
		log.Printf("Dealer choosen random %v %d", len(players), l)
		dealer = l
	} else {
		dealer = (previousDealer + 1) % len(players)
	}

	return &Hand{
		Players: players,
		//In:        players,
		Bank:      bank,
		Dealer:    dealer,
		HoleCards: make(map[string][2]models.Card, len(players)),
		InCount:   len(players),
		cardGen:   utils.NewCardSelector(),
		Blind:     10,
	}

}

func (h *Hand) WhileNotEnded(f func()) {
	if !h.Ended {
		f()
	}
}
