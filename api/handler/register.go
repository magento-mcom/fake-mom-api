package handler

import (
	"encoding/json"

	"github.com/magento-mcom/fake-mom-api/api"
	"github.com/pkg/errors"
	"fmt"
)

func NewRegisterHandler(registry api.Registry) api.Handler {
	return &registerHandler{
		registry: registry,
	}
}

type registerHandler struct {
	registry api.Registry
}

func (h *registerHandler) Handle(m *json.RawMessage) (interface{}, error) {
	i := api.Integration{}
	err := json.Unmarshal(*m, &i)

	if err != nil {
		return nil, err
	}

	if len(i.Url) <= 0 {
		return nil, errors.New("URL is empty")
	}

	h.registry.Add(i)

	fmt.Printf("Integration %v is registered with callback url: %v\n", i.ID, i.Url)

	return h.registry.GetAll(), nil
}
