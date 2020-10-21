package events

import (
	"errors"

	"github.com/JohnnyS318/go-poker/models"
	"github.com/mitchellh/mapstructure"
)

//FOLD descibes the action of a player quiting this hand
const FOLD = 1

//BET descibes the action of a player betting the same amount as the highes bet and therefore go along or callling the hand
const BET = 2

//RAISE raises sets the highest bet a certain amount
const RAISE = 3

//CHECK action pushes the action requirement to the next player
const CHECK = 4

//Action descibes a action a player can make one a normal hand stage
type Action struct {
	Action  int `json:"action" mapstructure:"action"`
	Payload int `json:"payload" mapstructure:"payload"`
}

func ToAction(raw *models.Event) (*Action, error) {

	if !ValidateEventName(PLAYER_ACTION, raw.Name) {
		return nil, errors.New(REQUIRED_EVENT_NAME_MISSING)
	}

	event := new(Action)
	err := mapstructure.Decode(raw.Data.(map[string]interface{}), event)
	return event, err
}

type WaitForActionEvent struct {
	Position int
	Preflop  bool
}

func NewWaitForActionEvent(position int, preflop bool) *models.Event {
	return models.NewEvent(WAIT_FOR_PLAYER_ACTION, &WaitForActionEvent{Position: position, Preflop: preflop})
}

type ActionProcessedEvent struct {
	Action   int                  `json:"action" mapstructure:"action"`
	Player   *models.PublicPlayer `json:"player" mapstructure:"player"`
	Position int                  `json:"positon" mapstructure:"position"`
	Amount   int                  `json:"amount" mapstructure:"amount"`
}

func NewActionProcessedEvent(action, amount, position int, player *models.Player) *models.Event {
	return models.NewEvent(ACTION_PROCESSED, &ActionProcessedEvent{
		Action:   action,
		Position: position,
		Amount:   amount,
		Player:   player.ToPublic(),
	})
}