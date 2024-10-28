package main

import (
	"sync"
)

const (
	SECRET_LENGTH = 1024
)

type broker interface {
	notify(message)

	subscribe(subscriber) string
	unsubscribe(string)

	done() <-chan struct{}
	close()
}

type brokerImpl struct {
	subList map[string]subscriber
	subLock sync.Mutex

	closed chan struct{}
}

func newBroker() *brokerImpl {
	return &brokerImpl{
		subList: make(map[string]subscriber),
	}
}

func (b *brokerImpl) close() {
	close(b.closed)
}

func (b *brokerImpl) done() <-chan struct{} {
	return b.closed
}

func (b *brokerImpl) subscribe(sub subscriber) string {
	id := genId()
	b.subLock.Lock()
	b.subList[id] = sub
	b.subLock.Unlock()
	return id
}

func (b *brokerImpl) unsubscribe(id string) {
	b.subLock.Lock()
	delete(b.subList, id)
	b.subLock.Unlock()
}

func (b *brokerImpl) notify(msg message) {
	b.subLock.Lock()
	for _, sub := range b.subList {
		logErr(sub.update(msg))
	}
	b.subLock.Unlock()
}
