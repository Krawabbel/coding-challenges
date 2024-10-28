package main

import (
	"io"
)

func handleSub(brok broker, writer io.Writer) error {
	sub := newSubscriber(writer)

	id := brok.subscribe(sub)
	defer brok.unsubscribe(id)

	<-sub.done()

	return nil
}

type subscriber interface {
	update(message) error
	done() <-chan struct{}
}

type subscriberImpl struct {
	writer io.Writer
	closed chan struct{}
}

func newSubscriber(writer io.Writer) *subscriberImpl {
	return &subscriberImpl{
		writer: writer,
		closed: make(chan struct{}),
	}
}

func (s *subscriberImpl) done() <-chan struct{} {
	return s.closed
}

func (s *subscriberImpl) update(msg message) error {

	if _, err := s.writer.Write([]byte(msg.string())); err != nil {
		close(s.closed)
		return err
	}

	return nil
}
