package events

import (
	"errors"
	"log"
	"reflect"

	"github.com/JohnnyS318/go-poker/models"
)

func NewWaitSmallBlindEvent(i int) *models.Event {
	return models.NewEvent(WAIT_FOR_SMALL_BLIND_SET, i)
}

func NewWaitBigBlindEvent(i int) *models.Event {
	return models.NewEvent(WAIT_FOR_BIG_BLIND_SET, i)
}

func ToBlindSet(e *models.Event) (int, error) {
	if !ValidateEventName(BLIND_SET, e.Name) {
		return -1, errors.New(REQUIRED_EVENT_NAME_MISSING)
	}

	log.Printf("data type %v", reflect.TypeOf(e.Data))

	r, ok := e.Data.(float64)
	if !ok {
		return -1, errors.New("Event data is not valid. Only valid bets are allowed")
	}
	return int(r), nil
}
