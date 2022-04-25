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

type MatchEnd struct {
	*elo.BaseEvent
	GameMode    string
	MapGroup    string
	MapFullName string
	MapName     string
	ScoreA      int
	ScoreB      int
	MatchStart  time.Time
	MatchEnd    time.Time
	Duration    time.Duration
	Completed   bool
}

func NewMatchEndEvent(b *elo.BaseEvent) (e *MatchEnd) {
	if gom := gameoverdrex.FindStringSubmatch(b.Message); gom != nil {
		mnpos := strings.LastIndex(gom[3], "/") + 1
		mn := gom[3][mnpos:]

		r := strings.NewReplacer(" ", "", "min", "m", "hour", "h", "sec", "s")
		dur := r.Replace(gom[6])
		d, err := time.ParseDuration(dur)
		if err != nil {
			log.Error("Could not parse game duration. Message: '%s'", b.Message)
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

		b.Server.CurrentMatch.GameMode = gom[1]
		b.Server.CurrentMatch.MapGroup = gom[2]
		b.Server.CurrentMatch.MapFullName = gom[3]
		b.Server.CurrentMatch.MapName = mn
		b.Server.CurrentMatch.ScoreA = scorea
		b.Server.CurrentMatch.ScoreB = scoreb
		b.Server.CurrentMatch.Completed = true

		b.Server.CurrentMatch.End = b.Time
		b.Server.CurrentMatch.Duration = d

		e = &MatchEnd{GameMode: b.Server.CurrentMatch.GameMode, MapGroup: b.Server.CurrentMatch.MapGroup,
			MapFullName: b.Server.CurrentMatch.MapFullName, MapName: b.Server.CurrentMatch.MapName,
			ScoreA: b.Server.CurrentMatch.ScoreA, ScoreB: b.Server.CurrentMatch.ScoreB,
			MatchEnd: b.Server.CurrentMatch.End, Duration: b.Server.CurrentMatch.Duration,
			Completed: b.Server.CurrentMatch.Completed,
			BaseEvent: b}
		log.Info("Created event: %#v", e)
	}
	return e
}

//Game Over: casual 2187570436 workshop/2209334999/de_elysion score 2:8 after 9 min
//Game Over: casual 2187570436 de_crete score 2:8 after 8 min
//Game Over: competitive default de_rats_brb score 8:3 after 16 min
//Game Over: competitive default de_shortnuke score 1:8 after 5 min
//Game Over: casual 2187570436 workshop/125444404/cs_office score 8:0 after 9 min
//Game Over: casual 2187570436 de_grind score 8:7 after 20 min
//Game Over: casual 2187570436 de_mocha score 5:8 after 16 min
var gameoverdrex = regexp.MustCompile(`Game Over: (.+) (.+) (.+) score (.+):(.+) after (.+)`)

func (e *MatchEnd) String() string {
	return fmt.Sprintf("Gameover %d:%d", e.ScoreA, e.ScoreB)
}
