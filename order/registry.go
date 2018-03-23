package order

type Registry interface {
	Save(id string)
	Exists(id string) bool
}

func NewOrderRegistry() Registry {
	return &registry{
		orders: map[string]bool{},
	}
}

type registry struct {
	orders map[string]bool
}

func (r *registry) Save(id string) {
	r.orders[id] = true
}

func (r *registry) Exists(id string) bool {
	_, ok := r.orders[id]

	return ok
}
