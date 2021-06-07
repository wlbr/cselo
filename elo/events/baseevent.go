package events

import (
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

type Event interface {
	String() string
}

type BaseEvent struct {
	Server  *elo.Server
	Time    time.Time
	Message string
}

func NewBaseEvent(server *elo.Server, timestamp time.Time, message string) *BaseEvent {
	e := &BaseEvent{Server: server, Time: timestamp, Message: message}
	log.Info("Created event: %+v", e)
	return e
}
