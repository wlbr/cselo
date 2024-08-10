package events

import (
	"fmt"
	"regexp"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type ServerHibernation struct {
	*elo.BaseEvent
}

func NewServerHibernationEvent(b *elo.BaseEvent) (e *ServerHibernation) {
	if hibrex.Match([]byte(b.Message)) {
		e = &ServerHibernation{BaseEvent: b}
		log.Info("Created event: %#v", e)
	}
	return e
}

// L 03/17/2022 - 22:07:11: "GOTV<42><BOT><Unassigned>" disconnected (reason "Punting bot, server is hibernating")
// "GOTV<42><BOT><Unassigned>" disconnected (reason "Punting bot, server is hibernating")
var hibrex = regexp.MustCompile(`"?GOTV<\d+><BOT><Unassigned>"? disconnected \(reason "?Punting bot, server is hibernating"?\)`)

func (e *ServerHibernation) String() string {
	return fmt.Sprintf("Hibernation of server %s at %v", e.Server.IP, e.Time)
}
