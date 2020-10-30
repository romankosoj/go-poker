package events

import (
	"errors"

	"github.com/JohnnyS318/go-poker/models"
	"github.com/mitchellh/mapstructure"
)

type JoinEvent struct {
	Username string `json:"username" mapstructure:"username"`
	ID       string `json:"id" mapstructure:"id"`
	BuyIn    int    `json:"buyin" mapstructure:"buyin"`
	LobbyID  string `json:"lobbyID" mapstructure:"lobbyId"`
}

func ToJoinEvent(raw *models.Event) (*JoinEvent, error) {

	if !ValidateEventName(JOIN, raw.Name) {
		return nil, errors.New(REQUIRED_EVENT_NAME_MISSING)
	}

	event := new(JoinEvent)
	err := mapstructure.Decode(raw.Data.(map[string]interface{}), event)
	return event, err
}

type JoinSuccess struct {
	LobbyID     string          `json:"lobbyId" mapstructure:"lobbyId"`
	Players     []models.Player `json:"players" mapstructure:"players"`
	GameStarted bool            `json:"gameStarted" mapstructure:"gameStarted"`
	MaxBuyIn    int             `json:"maxBuyIn"`
	MinBuyIn    int             `json:"minBuyIn"`
	BigBlind    int             `json:"bigBlind"`
	Position    int             `json:"position" mapstructure:"position"`
}

func NewJoinSuccessEvent(lobbyId string, players []models.Player, gameStarted bool, position, maxBuyIn, minBuyIn, bigBlind int) *models.Event {
	return models.NewEvent(JOIN_SUCCESS, &JoinSuccess{
		LobbyID:     lobbyId,
		Players:     players,
		GameStarted: gameStarted,
		MaxBuyIn:    maxBuyIn,
		MinBuyIn:    minBuyIn,
		BigBlind:    bigBlind,
		Position:    position,
	})
}

type PlayerLeavesEvent struct {
	Player *models.PublicPlayer `json:"player"`
	Index  int                  `json:"index"`
}

func NewPlayerLeavesEvent(player *models.Player, i int) *PlayerLeavesEvent {
	return &PlayerLeavesEvent{
		Player: player.ToPublic(),
		Index:  i,
	}
}
