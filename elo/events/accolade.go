package events

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type Accolade struct {
	BaseEvent
	Subject  *elo.Player
	Type     string
	Position int
	Value    float64
	Score    float64
}

// 3k  4k  5k  adr  assists  bombcarrierkills  burndamage  cashspent
// chickenskilled  clutchkills damage  deaths  dinks  enemiesflashed  firstkills
// gimme_01  gimme_02  gimme_03  gimme_04  gimme_05  gimme_06  headshotkills
// hsp  killreward  kills  killswhileblind  livetime  loudest  mvps
// nopurchasewins  objective  pistolkills  quietest  roundssurvived
// score  sniperkills  utilitydamage

//ACCOLADE, FINAL: {mvps},	Jacky<8>,	VALUE: 3.000000,	POS: 1,	SCORE: 24.000000
var accoladerexrex = regexp.MustCompile(`^ACCOLADE,\s+FINAL:\s+{(.+)},\s+(.+)<(.+)>,\s+VALUE:\s+(.+),\s+POS:\s+(.+),\s+SCORE:\s+(.+)$`)

func NewAccoladeEvent(server *elo.Server, t time.Time, message string) (a *Accolade) {
	if sm := accoladerexrex.FindStringSubmatch(message); sm != nil {
		pos, err := strconv.Atoi(sm[5])
		if err != nil {
			log.Error("Could not read position in accolade. %v   message: %s", err, message)
		} else {
			val, err := strconv.ParseFloat(sm[4], 64)
			if err != nil {
				log.Error("Could not read value in accolade. %v   message: %s", err, message)
			} else {
				sco, err := strconv.ParseFloat(sm[6], 64)
				if err != nil {
					log.Error("Could not read score in accolade. %v   message: %s", err, message)
				} else {
					p, err := elo.GetPlayerByName(sm[2])
					if err != nil {
						log.Info("Cannot identify player mentioned in accolade: %s   message: %s", sm[2], message)
					} else {
						a = &Accolade{Subject: p, Type: sm[1], Value: val, Position: pos, Score: sco,
							BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
						log.Info("Created event: %+v", a)
					}
				}
			}
		}
	}
	return a
}

func (e *Accolade) String() string {
	return fmt.Sprintf("Accolade: %s - %s", e.Subject, e.Type)
}
