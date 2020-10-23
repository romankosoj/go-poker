package hand

import "github.com/JohnnyS318/go-poker/models"

func checkIfEmpty(blocking [10]*models.Player) bool {
	for i := 0; i < 10; i++ {
		if blocking[i] != nil {
			return false
		}
	}
	return true
}

func removeBlocking(blocking [10]*models.Player, i int) {
	blocking[i] = nil
}

func addBlocking(blocking [10]*models.Player, i int, player *models.Player) error {
	isOn := false
	for _, n := range blocking {
		if n.ID == player.ID {
			isOn = true
		}
	}
	if !isOn {
		blocking[i] = player
	}

	return nil
}

func addAllButThisBlockgin(blocking [10]*models.Player, players []models.Player, player *models.Player) {
	for i, n := range players {
		if n.ID != player.ID && n.Active {
			blocking[i] = &n
		}
	}
}
