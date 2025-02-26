package util

import (
	"slices"
	"sync"
)

type MessageQueue struct {
	mutex    sync.Mutex
	channels map[string][]chan any
	buffer   int
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{
		mutex:    sync.Mutex{},
		channels: make(map[string][]chan any),
		buffer:   100,
	}
}

func (mq *MessageQueue) Subscribe(topic string) (<-chan any, func()) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	channel := make(chan any, mq.buffer)
	mq.channels[topic] = append(mq.channels[topic], channel)

	cancel := func() {
		mq.mutex.Lock()
		defer mq.mutex.Unlock()

		for i, c := range mq.channels[topic] {
			if channel == c {
				mq.channels[topic] = slices.Delete(mq.channels[topic], i, i+1)
				close(channel)
				return
			}
		}
	}

	return channel, cancel
}

func (mq *MessageQueue) Publish(topic string, payload any) {
	mq.mutex.Lock()
	defer mq.mutex.Unlock()

	for _, channel := range mq.channels[topic] {
		if len(channel) < mq.buffer {
			channel <- payload
		}
	}
}
