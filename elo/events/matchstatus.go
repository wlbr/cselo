package events

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type MatchStatus struct {
	*elo.BaseEvent
	MapFullName string
	MapName     string
	ScoreA      int
	ScoreB      int
	Rounds      int
}

func NewMatchStatusEvent(b *elo.BaseEvent) (e *MatchStatus) {
	if gom := matchstatusrex.FindStringSubmatch(b.Message); gom != nil {

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

		mnpos := strings.LastIndex(gom[3], "/") + 1
		mn := gom[3][mnpos:]

		//b.Server.CurrentMatch().SetStatusAttributes(gom[3], mn, scorea, scoreb, rounds)

		e = &MatchStatus{MapFullName: gom[3], MapName: mn, ScoreA: scorea, ScoreB: scoreb, Rounds: rounds, BaseEvent: b}
		log.Info("Created event: %#v", e)
	}
	return e
}

// CS2
//
// 11/04/2023 - 16:07:05.705 - MatchStatus: Score: 3:8 on map "de_inferno" RoundsPlayed: 11
//
// CSGO
// L 03/13/2022 - 14:09:17: MatchStatus: Score: 1:3 on map "workshop/2209334999/de_elysion" RoundsPlayed: 4
// MatchStatus: Score: 3:8 on map "de_inferno" RoundsPlayed: 11
// MatchStatus: Score: 2:0 on map "de_crete" RoundsPlayed: 2
// MatchStatus: Score: 4:3 on map "cs_italy" RoundsPlayed: 7
var matchstatusrex = regexp.MustCompile(`MatchStatus: Score: (\d+):(\d+) on map "?(.+)" RoundsPlayed: (\d+)`)

func (e *MatchStatus) String() string {
	return fmt.Sprintf("Matchstatus %d:%d", e.ScoreA, e.ScoreB)
}
