package events

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type MatchStart struct {
	BaseEvent
	MapFullName string
	MapName     string
}

//World triggered "Match_Start" on "de_lake"
//World triggered "Match_Start" on "workshop/123518981/de_favela"
var matchstartrex = regexp.MustCompile(`World triggered "Match_Start" on (.+)`)

func NewMatchStartEvent(server *elo.Server, t time.Time, message string) (e *MatchStart) {
	if sm := matchstartrex.FindStringSubmatch(message); sm != nil {
		mnpos := strings.LastIndex(sm[1], "/") + 1
		mn := sm[1][mnpos:]
		m := &elo.Match{MapFullName: sm[1], MapName: mn, Start: t, Server: server}
		server.CurrentMatch = m
		e = &MatchStart{MapFullName: m.MapFullName, MapName: m.MapName, BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *MatchStart) String() string {
	return fmt.Sprintf("Match start at %s", e.Time.Format(time.RFC822Z))
}
