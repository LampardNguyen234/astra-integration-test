package message_queue

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"
	"reflect"
	"sync"
)

type MessageQueue interface {
	QueueWriter
	QueueReader

	// Reset resets the queue.
	Reset() error
}

type QueueWriter interface {
	// IsFull checks if the MQ is full.
	IsFull() bool

	// CanAdd checks if the queue can add more queues.
	CanAdd(int) bool

	// Enqueue adds an element to the queue.
	Enqueue(interface{}) error

	// EnqueueBatch adds a number of queues to the end of the queue.
	EnqueueBatch([]interface{}) error
}

type QueueReader interface {
	// IsEmpty checks if the queue is empty.
	IsEmpty() bool

	// Dequeue gets the first element of the queue.
	Dequeue() interface{}

	// DequeueTo gets the first element of the queue and parses it to the given `ret` argument.
	DequeueTo(interface{}) error

	// MustDequeueTo gets the first element of the queue and parses it to the given `ret` argument.
	// If errors occur, it will panic. It should be used in accordance with IsEmpty method to avoid panicking.
	MustDequeueTo(interface{})
}

const (
	defaultQueueSize = uint64(math.MaxInt32)
)

// SimpleMQ implements a simple queue for handling messages.
type SimpleMQ struct {
	// size is the maximum number of elements.
	size uint64

	// elements is the
	elements []interface{}

	mtx *sync.Mutex

	elementType reflect.Type
}

// NewSimpleQueue creates a new SimpleMQ.
// If sizes is empty, defaultQueueSize will be used.
func NewSimpleQueue(sizes ...uint64) *SimpleMQ {
	size := defaultQueueSize
	if len(sizes) != 0 {
		size = sizes[0]
	}
	return &SimpleMQ{
		size:     size,
		elements: make([]interface{}, 0),
		mtx:      new(sync.Mutex),
	}
}

// IsFull checks if the queue is full.
func (q *SimpleMQ) IsFull() bool {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return len(q.elements) >= int(q.size)
}

// IsEmpty checks if the queue is empty.
func (q *SimpleMQ) IsEmpty() bool {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return len(q.elements) == 0
}

// CanAdd checks if the queue can add more queues.
func (q *SimpleMQ) CanAdd(size int) bool {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	return q.size-uint64(len(q.elements)) >= uint64(size)
}

// Enqueue adds an element to the end of the queue.
func (q *SimpleMQ) Enqueue(elem interface{}) error {
	if !q.CanAdd(1) {
		return ErrMQFull
	}
	t := reflect.TypeOf(elem)
	if q.elementType == nil {
		q.elementType = t
	}
	if t != q.elementType {
		return fmt.Errorf("expect to `%v` got `%v`", q.elementType, t)
	}

	q.mtx.Lock()
	defer q.mtx.Unlock()

	q.elements = append(q.elements, elem)
	return nil
}

// EnqueueBatch adds a number of queues to the end of the queue.
func (q *SimpleMQ) EnqueueBatch(elems []interface{}) error {
	if !q.CanAdd(len(elems)) {
		return ErrMQFull
	}

	q.mtx.Lock()
	defer q.mtx.Unlock()

	q.elements = append(q.elements, elems...)

	return nil
}

// Dequeue gets the first element of the queue.
func (q *SimpleMQ) Dequeue() interface{} {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if len(q.elements) == 0 {
		return nil
	}

	ret := q.elements[0]
	q.elements = q.elements[1:]

	return ret
}

// DequeueTo gets the first element of the queue and parses it to the given `ret` argument.
func (q *SimpleMQ) DequeueTo(ret interface{}) error {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	if len(q.elements) == 0 {
		return ErrMQEmpty
	}

	tmpRet := q.elements[0]
	jsb, _ := json.Marshal(tmpRet)
	err := json.Unmarshal(jsb, ret)
	if err != nil {
		return errors.Wrapf(ErrMQ, "failed to parse result: %v", err)
	}

	q.elements = q.elements[1:]

	return nil
}

// MustDequeueTo gets the first element of the queue and parses it to the given `ret` argument.
// If errors occur, it will panic. It should be used in accordance with IsEmpty method to avoid panicking.
func (q *SimpleMQ) MustDequeueTo(ret interface{}) {
	err := q.DequeueTo(ret)
	if err != nil {
		panic(err)
	}
}

func (q *SimpleMQ) Reset() error {
	q.mtx.Lock()
	defer q.mtx.Unlock()

	q.elements = make([]interface{}, 0)

	return nil
}
