package message_queue

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleQueue_Enqueue(t *testing.T) {
	q := NewSimpleQueue(defaultQueueSize)

	assert.NoError(t, q.Enqueue(100))
	assert.Error(t, q.Enqueue("Hello"))
}

func TestSimpleQueue_Dequeue(t *testing.T) {
	q := NewSimpleQueue(defaultQueueSize)

	err := q.Enqueue(1000)
	assert.NoError(t, err)

	m := q.Dequeue()
	fmt.Println(m, q.elements)
}

func TestSimpleQueue_DequeueTo(t *testing.T) {
	q := NewSimpleQueue(defaultQueueSize)

	valToAdd := 100
	err := q.Enqueue(valToAdd)
	assert.NoError(t, err)

	var ret int
	err = q.DequeueTo(&ret)
	assert.NoError(t, err)
	assert.Equal(t, valToAdd, ret)
	fmt.Println(ret)

	err = q.Enqueue(valToAdd)
	assert.NoError(t, err)

	var tmpRet string
	err = q.DequeueTo(&tmpRet)
	fmt.Println(err)
	assert.Error(t, err)
}
