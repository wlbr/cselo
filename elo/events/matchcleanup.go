package events

import (
	"fmt"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type MatchCleanUp struct {
	*elo.BaseEvent
	Match *elo.Match
}

func NewMatchCleanUpEvent(b *elo.BaseEvent, match *elo.Match) (e *MatchCleanUp) {
	e = &MatchCleanUp{Match: match, BaseEvent: b}
	log.Info("Created event: %#v", e)

	return e
}

func (e *MatchCleanUp) String() string {
	return fmt.Sprintf("MatchCleanUp match %d", e.Match.ID())
}
