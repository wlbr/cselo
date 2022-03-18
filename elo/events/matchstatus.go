package events

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type MatchStatus struct {
	BaseEvent
	MapFullName string
	MapName     string
	ScoreA      int
	ScoreB      int
	Rounds      int
}

func NewMatchStatusEvent(server *elo.Server, t time.Time, message string) (e *MatchStatus) {
	if gom := matchstatusrex.FindStringSubmatch(message); gom != nil {

		scorea, err := strconv.Atoi(gom[1])
		if err != nil {
			log.Error("Cannot read score of team A (got %s). %v", gom[1], err)
		}
		scoreb, err := strconv.Atoi(gom[2])
		if err != nil {
			log.Error("Cannot read score of team B (got %s). %v", gom[2], err)
		}
		rounds, err := strconv.Atoi(gom[4])
		if err != nil {
			log.Error("Cannot read rounds in match status (got %s). %v", gom[4], err)
		}
		server.CurrentMatch.MapFullName = gom[3]

		mnpos := strings.LastIndex(gom[3], "/") + 1
		mn := gom[3][mnpos:]
		server.CurrentMatch.MapName = mn

		server.CurrentMatch.ScoreA = scorea
		server.CurrentMatch.ScoreB = scoreb

		server.CurrentMatch.Rounds = rounds

		e = &MatchStatus{MapFullName: server.CurrentMatch.MapFullName, MapName: server.CurrentMatch.MapName,
			ScoreA: server.CurrentMatch.ScoreA, ScoreB: server.CurrentMatch.ScoreB,
			Rounds:    server.CurrentMatch.Rounds,
			BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %#v", e)
	}
	return e
}

//L 03/13/2022 - 14:09:17: MatchStatus: Score: 1:3 on map "workshop/2209334999/de_elysion" RoundsPlayed: 4
//MatchStatus: Score: 2:0 on map "de_crete" RoundsPlayed: 2
//MatchStatus: Score: 4:3 on map "cs_italy" RoundsPlayed: 7
var matchstatusrex = regexp.MustCompile(`MatchStatus: Score: (\d+):(\d+) on map \"(.+)\" RoundsPlayed: (\d+)`)

func (e *MatchStatus) String() string {
	return fmt.Sprintf("Matchstatus %d:%d", e.ScoreA, e.ScoreB)
}
