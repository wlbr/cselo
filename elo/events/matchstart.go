package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

type MatchStart struct {
	BaseEvent
}

//World triggered "Match_Start" on "de_lake"
var matchstartrex = regexp.MustCompile(`^World triggered "Match_Start".+$`)

func NewMatchStartEvent(server *elo.Server, t time.Time, message string) (e *MatchStart) {
	if sm := matchstartrex.FindStringSubmatch(message); sm != nil {
		e = &MatchStart{BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *MatchStart) String() string {
	return fmt.Sprintf("Match start at %s", e.Time.Format(time.RFC822Z))
}
