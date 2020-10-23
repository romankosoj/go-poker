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
}

func ToJoinEvent(raw *models.Event) (*JoinEvent, error) {

	if !ValidateEventName(JOIN, raw.Name) {
		return nil, errors.New(REQUIRED_EVENT_NAME_MISSING)
	}

	event := new(JoinEvent)
	err := mapstructure.Decode(raw.Data.(map[string]interface{}), event)
	return event, err
}

type PlayerLeavesEvent struct {
	Player *models.PublicPlayer `json:"player"`
	Index  int                  `index`
}

func NewPlayerLeavesEvent(player *models.Player, i int) *PlayerLeavesEvent {
	return &PlayerLeavesEvent{
		Player: player.ToPublic(),
		Index:  i,
	}
}
