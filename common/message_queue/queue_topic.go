package message_queue

import (
	"github.com/pkg/errors"
	"sync"
)

type TopicQueue interface {
	TopicQueueWriter
	TopicQueueReader

	// AddTopic adds a new topic queue.
	AddTopic(string, MessageQueue)

	// HasTopic checks if the TopicQueue has the given topic.
	HasTopic(string) bool
}

type TopicQueueWriter interface {
	// CanAdd checks if the topic queue can add more queues.
	CanAdd(string, int) bool

	// Enqueue adds an element to the topic queue.
	Enqueue(string, interface{}) error

	// EnqueueBatch adds a number of queues to the end of the topic queue.
	EnqueueBatch(string, []interface{}) error
}

type TopicQueueReader interface {
	// IsEmpty checks if the topic queue is empty.
	IsEmpty(string) bool

	// Dequeue gets the first element of the topic queue.
	Dequeue(string) interface{}

	// DequeueTo gets the first element of the topic queue and parses it to the given `ret` argument.
	DequeueTo(string, interface{}) error

	// MustDequeueTo gets the first element of the topic queue and parses it to the given `ret` argument.
	// If errors occur, it will panic. It should be used in accordance with IsEmpty method to avoid panicking.
	MustDequeueTo(string, interface{})
}

// SimpleTopicQueue implements a simple MQ management tool.
type SimpleTopicQueue struct {
	// queues is the
	queues map[string]MessageQueue

	mtx *sync.Mutex
}

// NewTopicQueue creates a new SimpleTopicQueue.
func NewTopicQueue(topics ...string) *SimpleTopicQueue {
	queues := make(map[string]MessageQueue)
	for _, topic := range topics {
		queues[topic] = NewSimpleQueue()
	}
	return &SimpleTopicQueue{
		queues: queues,
		mtx:    new(sync.Mutex),
	}
}

func (m *SimpleTopicQueue) AddTopic(topic string, queue MessageQueue) {
	if !m.HasTopic(topic) {
		m.queues[topic] = queue
	}
}

func (m *SimpleTopicQueue) HasTopic(topic string) bool {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return m.queues[topic] != nil
}

// IsEmpty checks if the topic queue is empty.
func (m *SimpleTopicQueue) IsEmpty(topic string) bool {
	if !m.HasTopic(topic) {
		return true
	}
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.queues[topic].IsEmpty()
}

// CanAdd checks if the queue can add more queues.
func (m *SimpleTopicQueue) CanAdd(topic string, size int) bool {
	if !m.HasTopic(topic) {
		return true
	}
	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.queues[topic].CanAdd(size)
}

// Enqueue adds an element to the end of the topic queue.
func (m *SimpleTopicQueue) Enqueue(topic string, elem interface{}) error {
	if !m.HasTopic(topic) {
		return errors.Wrapf(ErrTopicNotFound, "topic= %v", topic)
	}
	if !m.CanAdd(topic, 1) {
		return ErrMQFull
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.queues[topic].Enqueue(elem)
}

// EnqueueBatch adds a number of queues to the end of the topic queue.
func (m *SimpleTopicQueue) EnqueueBatch(topic string, elems []interface{}) error {
	if !m.HasTopic(topic) {
		return errors.Wrapf(ErrTopicNotFound, "topic= %v", topic)
	}
	if !m.CanAdd(topic, len(elems)) {
		return ErrMQFull
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.queues[topic].EnqueueBatch(elems)
}

// Dequeue gets the first element of the topic queue.
func (m *SimpleTopicQueue) Dequeue(topic string) interface{} {
	if !m.HasTopic(topic) {
		return nil
	}

	m.mtx.Lock()
	defer m.mtx.Unlock()

	return m.queues[topic].Dequeue()
}

// DequeueTo gets the first element of the topic queue and parses it to the given `ret` argument.
func (m *SimpleTopicQueue) DequeueTo(topic string, ret interface{}) error {
	if !m.HasTopic(topic) {
		return errors.Wrapf(ErrTopicNotFound, "topic= %v", topic)
	}

	return m.queues[topic].DequeueTo(ret)
}

// MustDequeueTo gets the first element of the topic queue and parses it to the given `ret` argument.
// If errors occur, it will panic. It should be used in accordance with IsEmpty method to avoid panicking.
func (m *SimpleTopicQueue) MustDequeueTo(topic string, ret interface{}) {
	err := m.DequeueTo(topic, ret)
	if err != nil {
		panic(err)
	}
}
