package consumer

import (
	"github.com/magento-mcom/fake-mom-api/api"
	"sync"
	"time"
	"errors"
)

type ConsumerQueue struct{
	requests []api.Request
	lock sync.Mutex
}

func (q *ConsumerQueue) Pop () (api.Request, error) {
	defer q.lock.Unlock()

	q.lock.Lock()

	if len(q.requests) == 0 {
		return api.Request{}, errors.New("Empty queue")
	}

	element := q.requests[0]
	q.requests = q.requests[1:]

	return element, nil
}

func NewConsumerQueue() *ConsumerQueue {
	return &ConsumerQueue{
		requests: []api.Request{},
		lock: sync.Mutex{},
	}
}

func NewConsumer(queue *ConsumerQueue, publisher api.Publisher) *Consumer {
	return &Consumer{
		queue:     queue,
		publisher: publisher,
	}
}

func (q *ConsumerQueue) Push (request api.Request) {
	defer q.lock.Unlock()

	q.lock.Lock()

	q.requests = append(q.requests, request)
}

type Consumer struct {
	queue *ConsumerQueue
	publisher api.Publisher
}

func (c *Consumer) Run () {
	for {
		time.Sleep(10 * time.Second)
		request, err := c.queue.Pop()
		if err == nil {
			c.publisher.Publish(request)
		}
	}
}
