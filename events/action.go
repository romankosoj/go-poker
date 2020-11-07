package events

import (
	"errors"
	"log"

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
	Position        int  `json:"position" mapstructure:"position"`
	PossibleActions byte `json:"possibleActions" mapstructure:"possibleActions"`
}

// NewWaitForAction is an event that the server is waiting for an action from a given player. The possible actions range from 0001 = Fold | 0010=Bet | 0100=Raise | 1000=Check to 1111=All
func NewWaitForActionEvent(position int, possibleActions byte) *models.Event {
	return models.NewEvent(WAIT_FOR_PLAYER_ACTION, &WaitForActionEvent{Position: position, PossibleActions: possibleActions})
}

type ActionProcessedEvent struct {
	Action int `json:"action" mapstructure:"action"`
	//Player   *models.PublicPlayer `json:"player" mapstructure:"player"`
	Position    int `json:"position" mapstructure:"position"`
	Amount      int `json:"amount" mapstructure:"amount"`
	TotalAmount int `json:"totalAmount" mapstructure:"totalAmount"`
}

func NewActionProcessedEvent(action, amount, position, totalAmount int) *models.Event {
	log.Printf("Player [%v] total: %v", position, totalAmount)
	return models.NewEvent(ACTION_PROCESSED, &ActionProcessedEvent{
		Action:      action,
		Position:    position,
		Amount:      amount,
		TotalAmount: totalAmount,
		//Player:   player.ToPublic(),
	})
}
