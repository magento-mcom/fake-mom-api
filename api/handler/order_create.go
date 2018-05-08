package handler

import (
	"encoding/json"

	"github.com/magento-mcom/fake-mom-api/api"
	"github.com/magento-mcom/fake-mom-api/order"
	"github.com/satori/go.uuid"
	"github.com/magento-mcom/fake-mom-api/consumer"
	"time"
	"fmt"
)

type OrderStatus struct {
	Status string `yaml:"status"`
	Reason string `yaml:"reason"`
}

type CreateOrderLine struct {
	Id string          `json:"id"`
	LineNumber int     `json:"line_number"`
	ProductName string `json:"product_name"`
	ProductType string `json:"product_type"`
	Sku string         `json:"sku"`
}

type CreateOrder struct {
	Order struct {
		Id string
		Store string
		Lines []CreateOrderLine
	}
}

type CustomerShipmentDoneItem struct {
	OrderLineId string  `json:"order_line_id"`
	OrderLineNumber int `json:"order_line_number"`
	ItemType string     `json:"item_type"`
	Sku string          `json:"sku"`
	Name string         `json:"name"`
}

type CustomerShipmentDonePackage struct {
	ID string   `json:"id"`
	Items []int `json:"items"`
}

type CustomerShipmentDone struct {
	Shipment struct{
		ShipmentId string                      `json:"shimpent_id"`
		StoreId string                         `json:"store_id"`
		OrderId string                         `json:"order_id"`
		SourceId string                        `json:"source_id"`
		Method string                          `json:"method"`
		ShipmentDate string                    `json:"shipment_date"`
		Items []CustomerShipmentDoneItem       `json:"items"`
		Packages []CustomerShipmentDonePackage `json:"packages"`
		Address struct{
			Reference string                   `json:"reference"`
			AddressType string                 `json:"address_type"`
			FirstName string                   `json:"first_name"`
			LastName string                    `json:"last_name"`
			Address1 string                    `json:"address1"`
			Address2 string                    `json:"address2"`
			City string                        `json:"city"`
			State string                       `json:"state"`
			Zip string                         `json:"zip"`
			CountryCode string                 `json:"country_code"`
			Phone string                       `json:"phone"`
			Email string                       `json:"email"`
		}                                      `json:"address"`
	}                                          `json:"shipment"`
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
	h.sendCustomerShipmentDone(m)

	return nil, nil
}

func (h *createOrderHandler) sendOrderCreated(message *json.RawMessage) {
	h.publishMessage("magento.sales.order_management.created", message)
}

func (h *createOrderHandler) sendOrderUpdated(message *json.RawMessage) {
	jsonMap := make(map[string]interface{})
	json.Unmarshal(*message, &jsonMap)
	for _, s := range h.statusToExport {
		jsonMap["order"].(map[string]interface{})["status"] = s.Status
		jsonMap["order"].(map[string]interface{})["status_reason"] = s.Reason
		m, _ := json.Marshal(jsonMap)
		params := json.RawMessage(m)
		h.publishMessage("magento.sales.order_management.updated", &params)
	}
}

func (h *createOrderHandler) sendCustomerShipmentDone(m CreateOrder) {
	cst := convertCreateOrderToCustomerShipmentDone(m)

	marsh, err := json.Marshal(cst)
	if err != nil {
		fmt.Printf("Could not serialize CustomerShipmenDone to JSON")
	} else {
		message := json.RawMessage(marsh)
		h.publishMessage("magento.logistics.fulfillment_management.customer_shipment_done", &message)
	}
}

func convertCreateOrderToCustomerShipmentDone(m CreateOrder) CustomerShipmentDone {
	cst := CustomerShipmentDone{}
	LinesShipped := []int{}
	cst.Shipment.ShipmentId = m.Order.Id + "-001"
	cst.Shipment.Method = "STANDARD"
	cst.Shipment.OrderId = m.Order.Id
	cst.Shipment.StoreId = m.Order.Store
	cst.Shipment.SourceId = "WAREHOUSE_1"
	cst.Shipment.ShipmentDate = time.Now().Format("2006-01-02 15:04:05-07:00")
	for _, i := range m.Order.Lines {
		si := CustomerShipmentDoneItem{
			ItemType:        i.ProductType,
			Name:            i.ProductName,
			OrderLineId:     i.Id,
			OrderLineNumber: i.LineNumber,
			Sku: i.Sku,
		}
		cst.Shipment.Items = append(cst.Shipment.Items, si)
		LinesShipped = append(LinesShipped, i.LineNumber)
	}
	cst.Shipment.Packages = append(cst.Shipment.Packages, CustomerShipmentDonePackage{
		ID:    m.Order.Id + "-001",
		Items: LinesShipped,
	})

	return cst
}

func (h *createOrderHandler) publishMessage(method string, params *json.RawMessage) {
	id, _ := uuid.NewV4()
	req := api.Request{
		Params: params,
		Method: method,
		ID: id.String(),
		Client: "FAKE",
	}

	h.queue.Push(req)
}
