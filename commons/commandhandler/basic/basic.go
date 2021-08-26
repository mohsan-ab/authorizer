package basic

import (
	"errors"
	"reflect"

	"github.com/mohsanabbas/authorizer/commons/eventsource"
)

// ErrInvalidID missing initial event
var ErrInvalidID = errors.New("Invalid ID, initial event missing")

// Handler contains the info to manage commands
type Handler struct {
	repository     *eventsource.Repository
	aggregate      reflect.Type
	bucket, subset string
}

// NewCommandHandler return a handler
func NewCommandHandler(repository *eventsource.Repository, aggregate eventsource.AggregateHandler, bucket, subset string) eventsource.CommandHandle {
	return &Handler{
		repository: repository,
		aggregate:  reflect.TypeOf(aggregate).Elem(),
		bucket:     bucket,
		subset:     subset,
	}
}

// Handle a command
func (h *Handler) Handle(command eventsource.Command) error {
	var err error

	version := command.GetVersion()
	aggregate := reflect.New(h.aggregate).Interface().(eventsource.AggregateHandler)

	if version != 0 {
		if err = h.repository.Load(aggregate, command.GetAggregateID()); err != nil {
			return err
		}
	}

	if err = aggregate.HandleCommand(command); err != nil {
		return err
	}

	// if not contain a valid ID,  the initial event (some like createAggreagate event) is missing
	if aggregate.GetID() == "" {
		return ErrInvalidID
	}

	if err = h.repository.Save(aggregate, version); err != nil {
		return err
	}

	if err = h.repository.PublishEvents(aggregate, h.bucket, h.subset); err != nil {
		return err
	}

	return nil
}
