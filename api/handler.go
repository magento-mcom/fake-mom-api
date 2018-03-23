package api

import "encoding/json"

type Handler interface {
	Handle(message *json.RawMessage) (interface{}, error)
}
