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
	*elo.BaseEvent
	MapFullName string
	MapName     string
	Match       *elo.Match
}

// World triggered "Match_Start" on "de_lake"
// World triggered "Match_Start" on "workshop/123518981/de_favela"
var matchstartrex = regexp.MustCompile(`World triggered "?Match_Start"? on "?(.+)"`)

func NewMatchStartEvent(b *elo.BaseEvent) (e *MatchStart) {
	if sm := matchstartrex.FindStringSubmatch(b.Message); sm != nil {
		mfn := strings.ReplaceAll(sm[1], `"`, "")
		mnpos := strings.LastIndex(mfn, "/") + 1
		mn := mfn[mnpos:]
		m := &elo.Match{MapFullName: mfn, MapName: mn, Start: b.Time, Server: b.Server}
		b.Server.LastMatch = b.Server.CurrentMatch
		b.Server.CurrentMatch = m
		e = &MatchStart{MapFullName: m.MapFullName, MapName: m.MapName, BaseEvent: b, Match: m}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *MatchStart) String() string {
	return fmt.Sprintf("Match start at %s on map %s", e.Time.Format(time.RFC822Z), e.MapName)
}
