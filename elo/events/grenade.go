package events

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type Grenade struct {
	*elo.BaseEvent
	Subject     *elo.Player
	subjectTeam string
	position    string
	GrenadeType string
	flashentity int
}

//"Jagger<19><STEAM_1:0:681607><TERRORIST>" threw flashbang [808 -247 -627] flashbang entindex 150)
//"Jagger<19><STEAM_1:0:681607><CT>" threw hegrenade [1303 -496 -638]
//"KiF Charlies Silence<16><STEAM_1:0:710013><TERRORIST>" threw smokegrenade [906 -291 -638]
//"AHA<199><STEAM_1:1:689719><CT>" threw molotov [1020 -969 -766]
var grenadedrex = regexp.MustCompile(`"(.+)<(.+)><(.+)><(.+)>" threw (.+) \[.+\]((.+) entindex (.+)\))?`)

func NewGrenadeEvent(b *elo.BaseEvent) (e *Grenade) {
	if sm := grenadedrex.FindStringSubmatch(b.Message); sm != nil {
		gtype := sm[5]
		entity := 0
		if gtype == "flashbang" {
			var err error
			entity, err = strconv.Atoi(sm[8])
			if err != nil {
				log.Error("Entity number not a int. Error: '%s', message: '%s'", err, b.Message)
			}
		}
		e = &Grenade{Subject: elo.GetPlayer(sm[1], sm[3]), subjectTeam: sm[4],
			GrenadeType: gtype, position: sm[6], flashentity: entity,
			BaseEvent: b}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Grenade) String() string {
	return fmt.Sprintf("%s throws grenade %s", e.Subject, e.GrenadeType)
}
