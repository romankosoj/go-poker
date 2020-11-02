package utils

import "github.com/JohnnyS318/go-poker/models"

//SendToAll is a utitlity for sending an event (message) to an entire array (lobby) of players.
func SendToAll(players []models.Player, event *models.Event) {
	for i := range players {
		players[i].Out <- event.ToRaw()
	}
}

//SendToPlayerInList is a utility for sending an event (message) to a specific player in a slice
func SendToPlayerInList(players []models.Player, i int, event *models.Event) {
	players[i].Out <- event.ToRaw()
}

//SendToPlayer is a utility for sending an event (message) to a given player
func SendToPlayer(player *models.Player, event *models.Event) {
	player.Out <- event.ToRaw()
}
