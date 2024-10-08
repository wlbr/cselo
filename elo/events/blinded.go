package events

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

// "Das Schnitzel<214><BOT><CT>" blinded for 2.39 by "Tina<217><BOT><CT>" from flashbang entindex 156
// "Jagger<19><STEAM_1:0:681607><TERRORIST>" blinded for 4.18 by "Dackel<2><STEAM_1:0:1770206><CT>" from flashbang entindex 209
var blindedrex = regexp.MustCompile(`"?(.+)<(.+)><(.+)><(.+)>"? blinded for (.+) by "?(.+)<(.+)><(.+)><(.+)>"? from flashbang entindex (\d+).*`)

type Blinded struct {
	*elo.BaseEvent
	Subject     *elo.Player
	subjectTeam string
	Object      *elo.Player
	objectTeam  string
	flashentity int
	Duration    float64
}

func NewBlindedEvent(b *elo.BaseEvent) (e *Blinded) {
	if sm := blindedrex.FindStringSubmatch(b.Message); sm != nil {
		dur, err1 := strconv.ParseFloat(sm[5], 32)
		if err1 != nil {
			log.Error("Flash duration not a float. %s, message: \"%s\"", err1, b.Message)
		}
		ent, err2 := strconv.Atoi(sm[10])
		if err2 != nil {
			log.Error("Entity number not a int. %s, message: \"%s\" %v", err2, b.Message, ent)
		}
		if err1 == nil && err2 == nil {
			e = &Blinded{Subject: elo.GetPlayer(sm[6], sm[8]), subjectTeam: sm[9],
				Object: elo.GetPlayer(sm[1], sm[3]), objectTeam: sm[4],
				Duration: dur, flashentity: ent, BaseEvent: b}
		}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Blinded) String() string {
	s := fmt.Sprintf("Blinded %s ==> %s", e.Subject, e.Object)

	return s
}

func (e *Blinded) OwnTeamFlashed() bool {
	if e.subjectTeam == e.objectTeam {
		return true
	}
	return false
}

func (e *Blinded) EnemyFlashed() bool {
	if e.subjectTeam != e.objectTeam {
		return true
	}
	return false
}

func (e *Blinded) TeammateFlashed() bool {
	if e.subjectTeam == e.objectTeam && e.Subject.SteamID != e.Object.SteamID {
		return true
	}
	return false
}

func (e *Blinded) SelfFlashed() bool {
	if e.Subject.SteamID == e.Object.SteamID {
		return true
	}
	return false
}
