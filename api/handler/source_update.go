package handler

import (
	"encoding/json"

	"github.com/magento-mcom/fake-mom-api/api"
	"github.com/satori/go.uuid"
)

type SourceStockUpdate struct {
	Snapshot struct {
		SourceId string
		Mode     string
		Stock    []struct {
			Sku       string
			Quantity  int
			Unlimited bool
		}
	}
}

type AggregateStockUpdated struct {
	Snapshot struct {
		AggregateId string
		Mode        string
		Stock       []struct {
			Sku       string
			Quantity  int
			Unlimited bool
		}
	}
}

func NewSourceUpdateHandler(publisher api.Publisher, aggregatesToExport []string) api.Handler {
	return &sourceUpdateHandler{
		publisher:          publisher,
		aggregatesToExport: aggregatesToExport,
	}
}

type sourceUpdateHandler struct {
	publisher          api.Publisher
	aggregatesToExport []string
}

func (h *sourceUpdateHandler) Handle(message *json.RawMessage) (interface{}, error) {
	ss := SourceStockUpdate{}
	json.Unmarshal(*message, ss)

	as := AggregateStockUpdated{}

	as.Snapshot.Mode = ss.Snapshot.Mode
	as.Snapshot.Stock = ss.Snapshot.Stock

	req := api.Request{
		Method: "magento.inventory.aggregate_stock_management.updated",
		Client: "FAKE",
	}
	for _, a := range h.aggregatesToExport {
		as.Snapshot.AggregateId = a
		b, _ := json.Marshal(as)
		m := json.RawMessage(b)
		req.Params = &m
		id, _ := uuid.NewV4()
		req.ID = id.String()
		h.publisher.Publish(req)
	}

	return nil, nil
}
