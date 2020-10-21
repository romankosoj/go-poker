package lobbies

func (l *LobbyManager) Search() []int {

	a := make([]int, 0)

	for i, n := range l.Lobbies {
		if n.HasCapacaty() {
			a = append(a, i)
		}
	}
	return a
}
