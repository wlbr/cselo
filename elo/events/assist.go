package events

import (
	"fmt"
	"regexp"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

//"Tina<217><BOT><CT>" assisted killing "Franzi<216><BOT><TERRORIST>"
//"Jagger<19><STEAM_1:0:681607><TERRORIST>" assisted killing "AHA<199><STEAM_1:1:689719><CT>"
var assistrex = regexp.MustCompile(`"(.+)<(.+)><(.+)><(.+)>" assisted killing "(.+)<(.+)><(.+)><(.+)>"`)

type Assist struct {
	*elo.BaseEvent
	Subject     *elo.Player
	subjectTeam string
	Object      *elo.Player
	objectTeam  string
}

func NewAssistEvent(b *elo.BaseEvent) (e *Assist) {
	if sm := assistrex.FindStringSubmatch(b.Message); sm != nil {
		e = &Assist{Subject: elo.GetPlayer(sm[1], sm[3]), subjectTeam: sm[4],
			Object: elo.GetPlayer(sm[5], sm[7]), objectTeam: sm[8],
			BaseEvent: b}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Assist) String() string {
	return fmt.Sprintf("Assist %s ==> %s", e.Subject, e.Object)
}
