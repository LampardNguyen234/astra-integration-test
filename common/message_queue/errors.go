package message_queue

import "fmt"

var (
	ErrMQ            = fmt.Errorf("queue error")
	ErrMQEmpty       = fmt.Errorf("queue empty")
	ErrMQFull        = fmt.Errorf("queue full")
	ErrTopicNotFound = fmt.Errorf("topic not found")
)
