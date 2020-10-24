package utils

import (
	"log"
	"math/rand"

	"github.com/JohnnyS318/go-poker/models"
)

//CardGenerator is a utility for randomly choose a multiple cards out of a 52 card deck
type CardGenerator struct {
	Cards         [52]models.Card
	SelectedCards int
}

//NewCardSelector creates a new card select
func NewCardSelector() *CardGenerator {
	s := &CardGenerator{SelectedCards: 0}
	s.Reset()
	return s
}

//Reset resets the card stack and therefore the probabilities
func (s *CardGenerator) Reset() {
	for i := 0; i < 48; i += 4 {
		for j := 0; j < 4; j++ {

			c := new(models.Card)
			c, err := models.NewCard(j, i/4)

			if err != nil {
				log.Printf("Something went wrong %v", err)

				// Card selection is a key feature and should not fail
				// => closing lobby
			}
			s.Cards[i+j] = *c
		}
	}
	s.SelectedCards = 0
}

//SelectRandom randomly selects a card from the stack and returns a copy of it
func (s *CardGenerator) SelectRandom() models.Card {
	i := rand.Intn(51 - s.SelectedCards)
	c := s.Cards[i]
	// if selected card is last. the last card has not to be swapped
	if i != 51-s.SelectedCards {
		s.Cards[51-s.SelectedCards], s.Cards[i] = c, s.Cards[51-s.SelectedCards]
		s.SelectedCards++
	}
	return c
}
