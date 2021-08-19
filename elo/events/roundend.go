package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type RoundEnd struct {
	BaseEvent
}

//World triggered "Round_End"
var roundendedrex = regexp.MustCompile(`^World triggered "Round_End"$`)

func NewRoundEndEvent(server *elo.Server, t time.Time, message string) (e *RoundEnd) {
	if sm := roundstartdrex.FindStringSubmatch(message); sm != nil {
		e = &RoundEnd{BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		server.LastPlanter = nil
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *RoundEnd) String() string {
	return fmt.Sprintf("Round end at %s", e.Time.Format(time.RFC822Z))
}
