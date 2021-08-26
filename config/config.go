package config

import (
	"github.com/mohsanabbas/authorizer/commons/commandbus/async"
	"github.com/mohsanabbas/authorizer/commons/eventbus/nats"
	"github.com/mohsanabbas/authorizer/commons/eventsource"
	"github.com/mohsanabbas/authorizer/commons/eventstore/mongo"
)

// EventBus returns an eventsource.EventBus impl
type EventBus func() (eventsource.EventBus, error)

// EventStore returns an eventsource.EventStore impl
type EventStore func() (eventsource.EventStore, error)

// CommandBus returns an eventsource.CommandBus
type CommandBus func(register eventsource.CommandHandlerRegister) (eventsource.CommandBus, error)

// CommandConfig should connect internally commands with an aggregate
type CommandConfig func(repository *eventsource.Repository, register *eventsource.CommandRegister)

// commandHandler is the signature used by command handlers constructor
type commandHandler func(repository *eventsource.Repository, aggregate eventsource.AggregateHandler, bucket, subset string) eventsource.CommandHandle

// WireCommands acts as a wired between aggregate, register and commands
func WireCommands(aggregate eventsource.AggregateHandler, handler commandHandler, bucket, subset string, commands ...interface{}) CommandConfig {
	return func(repository *eventsource.Repository, register *eventsource.CommandRegister) {
		h := handler(repository, aggregate, bucket, subset)
		for _, command := range commands {
			register.Add(command, h)
		}
	}
}

// NewClient returns a command bus properly configured
func NewClient(es EventStore, eb EventBus, cb CommandBus, cmdConfigs ...CommandConfig) (eventsource.CommandBus, error) {
	store, err := es()
	if err != nil {
		return nil, err
	}

	bus, err := eb()
	if err != nil {
		return nil, err
	}

	repository := eventsource.NewRepository(store, bus)
	register := eventsource.NewCommandRegister()

	for _, conf := range cmdConfigs {
		conf(repository, register)
	}
	return cb(register)
}

// Nats generates a Nats implementation of EventBus
func Nats(urls string, useTLS bool) EventBus {
	return func() (eventsource.EventBus, error) {
		return nats.NewClient(urls, useTLS)
	}
}

// Mongo generates a MongoDB implementation of EventStore
func Mongo(host string, port int, db string) EventStore {
	return func() (eventsource.EventStore, error) {
		return mongo.NewClient(host, port, db)
	}
}

// AsyncCommandBus generates a CommandBus
func AsyncCommandBus(workers int) CommandBus {
	return func(register eventsource.CommandHandlerRegister) (eventsource.CommandBus, error) {
		return async.NewBus(register, workers), nil
	}
}
