package elo

import (
	"time"

	"github.com/wlbr/commons/log"
)

type Event interface {
	String() string
}

type BaseEvent struct {
	Server  *Server
	Time    time.Time
	Message string
}

func NewBaseEvent(server *Server, timestamp time.Time, message string) *BaseEvent {
	e := &BaseEvent{Server: server, Time: timestamp, Message: message}
	log.Info("Created event: %+v", e)
	return e
}
