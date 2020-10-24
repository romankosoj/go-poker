package hand

import (
	"log"
	"sort"

	"github.com/JohnnyS318/go-poker/models"
)

func checkIfEmpty(blocking []int) bool {
	return len(blocking) <= 0
}

func removeBlocking(blocking []int, i int) []int {
	b := append(blocking[:i], blocking[i+1:]...)
	log.Printf("Removed player [%v] now blocking: %v", i, b)
	return b
}

func addBlocking(blocking []int, k int) error {
	isOn := false
	for _, n := range blocking {
		if n == k {
			isOn = true
		}
	}
	if !isOn {
		blocking = append(blocking, k)
	}
	sort.Slice(blocking, func(i, j int) bool {
		return blocking[i] < blocking[j]
	})

	return nil
}

func addAllButThisBlockgin(blocking []int, players []models.Player, k int) []int {
	blocking = nil
	for i := range players {
		if i != k && players[i].Active {
			blocking = append(blocking, i)
		}
	}
	return blocking
}
