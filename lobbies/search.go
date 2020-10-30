package lobbies

func (l *LobbyManager) Search() []string {

	a := make([]string, 0)

	for k, v := range l.Lobbies {
		if v.HasCapacaty() {
			a = append(a, k)
		}
	}
	return a
}
