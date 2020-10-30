package lobbies

import (
	"github.com/JohnnyS318/go-poker/models"
)

type PlayerQueue struct {
	Queue    []*models.Player
	Register map[string]string
}

func NewPlayerQueue() *PlayerQueue {
	return &PlayerQueue{
		Queue:    make([]*models.Player, 0),
		Register: make(map[string]string),
	}
}

func (p *PlayerQueue) Enqueue(player *models.Player, id string) {
	p.Queue = append(p.Queue, player)
	if id != "" {
		p.Register[player.ID] = id
	}
}

func (p *PlayerQueue) Dequeue() (*models.Player, string) {
	player := p.Queue[0]
	p.Queue[0] = nil

	p.Queue = p.Queue[1:]

	id, ok := p.Register[player.ID]

	if !ok {
		return player, ""
	}

	delete(p.Register, id)

	return player, id

}
