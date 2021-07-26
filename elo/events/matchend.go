package events

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

type MatchEnd struct {
	BaseEvent
	GameMode    string
	MapGroup    string
	MapFullName string
	MapName     string
	ScoreA      int64
	ScoreB      int64
	MatchStart  time.Time
	MatchEnd    time.Time
	Duration    time.Duration
}

func NewMatchEndEvent(server *elo.Server, t time.Time, message string) (e *MatchEnd) {
	if gom := gameoverdrex.FindStringSubmatch(message); gom != nil {
		mnpos := strings.LastIndex(gom[3], "/") + 1
		mn := gom[3][mnpos:]

		r := strings.NewReplacer(" ", "", "min", "m", "hour", "h", "sec", "s")
		dur := r.Replace(gom[6])
		d, err := time.ParseDuration(dur)
		if err != nil {
			log.Error("Could not parse game duration. Message: '%s'", message)
			d, _ = time.ParseDuration("0m")
		}
		scorea, err := strconv.Atoi(gom[4])
		if err != nil {
			log.Error("Cannot read score of team A (got %s). %v", gom[4], err)
		}
		scoreb, err := strconv.Atoi(gom[5])
		if err != nil {
			log.Error("Cannot read score of team B (got %s). %v", gom[5], err)
		}

		server.CurrentMatch.GameMode = gom[1]
		server.CurrentMatch.MapGroup = gom[2]
		server.CurrentMatch.MapFullName = gom[3]
		server.CurrentMatch.MapName = mn
		server.CurrentMatch.ScoreA = int64(scorea)
		server.CurrentMatch.ScoreB = int64(scoreb)
		server.CurrentMatch.End = t
		server.CurrentMatch.Duration = d

		e = &MatchEnd{GameMode: server.CurrentMatch.GameMode, MapGroup: server.CurrentMatch.MapGroup,
			MapFullName: server.CurrentMatch.MapFullName, MapName: server.CurrentMatch.MapName,
			ScoreA: server.CurrentMatch.ScoreA, ScoreB: server.CurrentMatch.ScoreB,
			MatchEnd: server.CurrentMatch.End, Duration: server.CurrentMatch.Duration,
			BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %#v", e)
	}
	return e
}

//Game Over: competitive default de_rats_brb score 8:3 after 16 min
//Game Over: competitive default de_shortnuke score 1:8 after 5 min
//Game Over: casual 2187570436 workshop/125444404/cs_office score 8:0 after 9 min
//Game Over: casual 2187570436 de_grind score 8:7 after 20 min
//Game Over: casual 2187570436 de_mocha score 5:8 after 16 min
var gameoverdrex = regexp.MustCompile(`Game Over: (.+) (.+) (.+) score (.+):(.+) after (.+)`)

func (e *MatchEnd) String() string {
	return fmt.Sprintf("Gameover %d:%d", e.ScoreA, e.ScoreB)
}
