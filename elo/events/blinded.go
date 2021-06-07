package events

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

//"Das Schnitzel<214><BOT><CT>" blinded for 2.39 by "Tina<217><BOT><CT>" from flashbang entindex 156
//"Jagger<19><STEAM_1:0:681607><TERRORIST>" blinded for 4.18 by "Dackel<2><STEAM_1:0:1770206><CT>" from flashbang entindex 209
var blindedrex = regexp.MustCompile(`^"(.+)<(.+)><(.+)><(.+)>" blinded for (.+) by "(.+)<(.+)><(.+)><(.+)>" from flashbang entindex (\d+).*$`)

type Blinded struct {
	BaseEvent
	Subject     *elo.Player
	subjectTeam string
	Object      *elo.Player
	objectTeam  string
	flashentity int
	Duration    float64
}

func NewBlindedEvent(server *elo.Server, t time.Time, message string) (e *Blinded) {
	if sm := blindedrex.FindStringSubmatch(message); sm != nil {
		dur, err1 := strconv.ParseFloat(sm[5], 32)
		if err1 != nil {
			log.Error("Flash duration not a float. %s, message: \"%s\"", err1, message)
		}
		ent, err2 := strconv.Atoi(sm[10])
		if err2 != nil {
			log.Error("Entity number not a int. %s, message: \"%s\" %v", err2, message, ent)
		}
		if err1 == nil && err2 == nil {
			e = &Blinded{Subject: &elo.Player{Name: sm[6], SteamID: sm[8]}, subjectTeam: sm[9],
				Object: &elo.Player{Name: sm[1], SteamID: sm[3]}, objectTeam: sm[4],
				Duration: dur, flashentity: ent,
				BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Blinded) String() string {
	s := fmt.Sprintf("Blinded %s ==> %s", e.Subject, e.Object)

	return s
}

func (e *Blinded) OwnTeam() bool {
	if e.subjectTeam == e.objectTeam {
		return true
	}
	return false
}
