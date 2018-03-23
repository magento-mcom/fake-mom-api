package handler

import (
	"encoding/json"

	"github.com/magento-mcom/fake-api/api"
	"github.com/satori/go.uuid"
)

type OrderStatus struct {
	Status string `yaml:"status"`
	Reason string `yaml:"reason"`
}

func NewCreateOrderHandler(publisher api.Publisher, statusToExport []OrderStatus) api.Handler {
	return &createOrderHandler{
		publisher:      publisher,
		statusToExport: statusToExport,
	}
}

type createOrderHandler struct {
	publisher      api.Publisher
	statusToExport []OrderStatus
}

func (h *createOrderHandler) Handle(message *json.RawMessage) (interface{}, error) {
	h.sendOrderCreated(message)
	h.sendOrderUpdated(message)

	return nil, nil
}

func (h *createOrderHandler) sendOrderCreated(message *json.RawMessage) {
	id, _ := uuid.NewV4()
	req := api.Request{
		Params: message,
		Method: "magento.sales.order_management.created",
		ID:     id.String(),
		Client: "FAKE",
	}
	h.publisher.Publish(req)
}

func (h *createOrderHandler) sendOrderUpdated(message *json.RawMessage) {
	jsonMap := make(map[string]interface{})
	json.Unmarshal(*message, &jsonMap)
	for _, s := range h.statusToExport {
		jsonMap["order"].(map[string]interface{})["status"] = s.Status
		jsonMap["order"].(map[string]interface{})["status_reason"] = s.Reason
		m, _ := json.Marshal(jsonMap)
		params := json.RawMessage(m)
		id, _ := uuid.NewV4()
		req := api.Request{
			Params: &params,
			Method: "magento.sales.order_management.updated",
			ID:     id.String(),
			Client: "FAKE",
		}

		h.publisher.Publish(req)
	}
}
