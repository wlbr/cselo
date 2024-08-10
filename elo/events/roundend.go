package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type RoundEnd struct {
	*elo.BaseEvent
}

// World triggered "Round_End"
var roundendedrex = regexp.MustCompile(`World triggered "?Round_End"?`)

func NewRoundEndEvent(b *elo.BaseEvent) (e *RoundEnd) {
	if sm := roundstartdrex.FindStringSubmatch(b.Message); sm != nil {
		e = &RoundEnd{BaseEvent: b}
		b.Server.LastPlanter = nil
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *RoundEnd) String() string {
	return fmt.Sprintf("Round end at %s", e.Time.Format(time.RFC822Z))
}
