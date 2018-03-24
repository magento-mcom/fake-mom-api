package handler

import (
	"encoding/json"

	"github.com/magento-mcom/fake-mom-api/api"
	"github.com/magento-mcom/fake-mom-api/order"
	"github.com/satori/go.uuid"
	"github.com/magento-mcom/fake-mom-api/consumer"
)

type OrderStatus struct {
	Status string `yaml:"status"`
	Reason string `yaml:"reason"`
}

type CreateOrder struct {
	Order struct {
		Id string
	}
}

func NewCreateOrderHandler(c *consumer.ConsumerQueue, statusToExport []OrderStatus, registry order.Registry) api.Handler {
	return &createOrderHandler{
		queue:          c,
		statusToExport: statusToExport,
		registry:       registry,
	}
}

type createOrderHandler struct {
	queue          *consumer.ConsumerQueue
	statusToExport []OrderStatus
	registry       order.Registry
}

func (h *createOrderHandler) Handle(message *json.RawMessage) (interface{}, error) {
	m := CreateOrder{}
	json.Unmarshal(*message, &m)
	h.registry.Save(m.Order.Id)
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
	h.queue.Push(req)
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

		h.queue.Push(req)
	}
}
