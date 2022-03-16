package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type RoundStart struct {
	BaseEvent
}

//World triggered "Round_Start"
var roundstartdrex = regexp.MustCompile(`World triggered "Round_Start"`)

func NewRoundStartEvent(server *elo.Server, t time.Time, message string) (e *RoundStart) {
	if sm := roundstartdrex.FindStringSubmatch(message); sm != nil {
		e = &RoundStart{BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		server.LastPlanter = nil
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *RoundStart) String() string {
	return fmt.Sprintf("Round start at %s", e.Time.Format(time.RFC822Z))
}
