package models

type PublicPlayer struct {
	Username string `json:"username" mapstructure:"username"`
	ID       string `json:"id" mapstructure:"id"`
}

func (player *Player) ToPublic() *PublicPlayer {
	return &PublicPlayer{
		ID:       player.ID,
		Username: player.Username,
	}
}
