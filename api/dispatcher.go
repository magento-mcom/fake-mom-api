package api

import (
	"github.com/pkg/errors"
	"fmt"
)

func NewDispatcher(handlers map[string]Handler) Dispatcher {
	return &dispatcher{
		handlers: handlers,
	}
}

type Dispatcher interface {
	Dispatch(message Request) (interface{}, error)
}

type dispatcher struct {
	handlers map[string]Handler
}

func (d *dispatcher) Dispatch(message Request) (interface{}, error) {
	h, ok := d.handlers[message.Method]
	if !ok {
		return nil, errors.Errorf("Handler not found for method %v.", message.Method)
	}

	fmt.Printf("Received message with method %v\n", message.Method)

	return h.Handle(message.Params)
}
