package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
	"time"
)

func NewPublisher(registry Registry) Publisher {
	return &publisher{
		registry: registry,
		client:   http.Client{},
	}
}

type Publisher interface {
	Publish(request Request)
}

type publisher struct {
	registry Registry
	client   http.Client
}

func (p *publisher) Publish(request Request) {
	integrations := p.registry.GetAll()

	b, _ := json.Marshal(request)
	buffer := bytes.NewBuffer(b)

	for _, i := range integrations {
		time.Sleep(10 * time.Second)
		fmt.Printf("Publishing message with method %v to integration %v\n", request.Method, i.ID)
		p.client.Post(i.Url, "application/json", buffer)
	}
}
