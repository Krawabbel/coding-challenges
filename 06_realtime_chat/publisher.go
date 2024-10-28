package main

import (
	"errors"
	"io"
)

func handlePub(brok broker, reader io.Reader) error {
	pub := newPublisher(brok, reader)
	if err := pub.publish(); err != nil {
		return err
	}
	return nil
}

type publisher interface {
	publish() error
}

type publisherImpl struct {
	brok   broker
	reader io.Reader
}

func newPublisher(broker broker, reader io.Reader) *publisherImpl {
	return &publisherImpl{
		brok:   broker,
		reader: reader,
	}
}

func (p *publisherImpl) query() (message, error) {
	buf := make([]byte, MESSAGE_LENGTH)
	n, err := p.reader.Read(buf)
	if err != nil {
		return nil, err
	}
	return newMessage(buf[:n]), nil
}

func (p *publisherImpl) publish() error {
	for {
		msg, err := p.query()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		p.brok.notify(msg)
	}
}
