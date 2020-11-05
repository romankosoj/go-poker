package models

type PublicPlayer struct {
	Username string `json:"username" mapstructure:"username"`
	ID       string `json:"id" mapstructure:"id"`
	BuyIn    int    `json:"buyIn" mapstructure:"buyIn"`
}

func (player *Player) ToPublic() *PublicPlayer {
	return &PublicPlayer{
		ID:       player.ID,
		Username: player.Username,
		BuyIn:    player.BuyIn,
	}
}
