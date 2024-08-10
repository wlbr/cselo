package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type RoundStart struct {
	*elo.BaseEvent
}

// World triggered "Round_Start"
var roundstartdrex = regexp.MustCompile(`World triggered "?Round_Start"?`)

func NewRoundStartEvent(b *elo.BaseEvent) (e *RoundStart) {
	if sm := roundstartdrex.FindStringSubmatch(b.Message); sm != nil {
		e = &RoundStart{BaseEvent: b}
		b.Server.LastPlanter = nil
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *RoundStart) String() string {
	return fmt.Sprintf("Round start at %s", e.Time.Format(time.RFC822Z))
}
