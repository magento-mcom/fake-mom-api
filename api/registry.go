package api

import "sync"

type Integration struct {
	ID  string
	Url string
}

func NewRegistry() Registry {
	return &registry{
		integrations: []Integration{},
	}
}

type Registry interface {
	Add(integration Integration)
	GetAll() []Integration
}

type registry struct {
	integrations []Integration
	lock         sync.Mutex
}

func (r *registry) Add(integration Integration) {
	defer r.lock.Unlock()
	r.lock.Lock()
	r.integrations = append(r.integrations, integration)
}

func (r *registry) GetAll() []Integration {
	return r.integrations
}
