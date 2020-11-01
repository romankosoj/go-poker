package events

import "github.com/JohnnyS318/go-poker/models"

type GameEndEvent struct {
	Winners []models.PublicPlayer `json:"winners"`
	Share   int                   `json:"share"`
}

func NewGameEndEvent(winners []models.PublicPlayer, share int) *models.Event {
	return models.NewEvent(GAME_END, &GameEndEvent{
		Winners: winners,
		Share:   share,
	})
}
