package events

import (
	"fmt"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type MatchCleanUp struct {
	BaseEvent
	Match *elo.Match
}

func NewMatchCleanUpEvent(server *elo.Server, t time.Time, message string, match *elo.Match) (e *MatchCleanUp) {
	e = &MatchCleanUp{Match: match, BaseEvent: BaseEvent{Time: time.Now(), Server: server, Message: message}}
	log.Info("Created event: %#v", e)

	return e
}

func (e *MatchCleanUp) String() string {
	return fmt.Sprintf("MatchCleanUp  match %d%d", e.Match.ID)
}
