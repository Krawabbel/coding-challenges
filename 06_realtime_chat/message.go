package main

import (
	"fmt"
	"time"
)

const (
	MESSAGE_LENGTH = 1024
)

type message interface {
	string() string
}

type messageImpl struct {
	Data []byte    `json:"data"`
	Time time.Time `json:"time"`
}

func newMessage(data []byte) *messageImpl {
	return &messageImpl{
		Data: data,
		Time: time.Now(),
	}
}

func (m *messageImpl) string() string {
	return fmt.Sprint(m.Data)
}
