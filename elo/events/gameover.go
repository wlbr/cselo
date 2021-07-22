package events

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

type GameOver struct {
	BaseEvent
	GameMode    string
	MapGroup    string
	MapFullName string
	MapName     string
	Score       string
	Duration    time.Duration
}

func NewGameOverEvent(server *elo.Server, t time.Time, message string) (e *GameOver) {
	if gom := gameoverdrex.FindStringSubmatch(message); gom != nil {
		mnpos := strings.LastIndex(gom[3], "/") + 1
		mn := gom[3][mnpos:]

		r := strings.NewReplacer(" ", "", "min", "m", "hour", "h", "sec", "s")
		dur := r.Replace(gom[5])
		d, err := time.ParseDuration(dur)
		if err != nil {
			log.Error("Could not parse game duration. Message: '%s'", message)
			d, _ = time.ParseDuration("0m")
		}

		e = &GameOver{GameMode: gom[1], MapGroup: gom[2], MapFullName: gom[3], MapName: mn, Score: gom[4], Duration: d,
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
var gameoverdrex = regexp.MustCompile(`Game Over: (.+) (.+) (.+) score (.+) after (.+)`)

func (e *GameOver) String() string {
	return fmt.Sprintf("Gameover %s", e.Score)
}
