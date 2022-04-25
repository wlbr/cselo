package events

import (
	"fmt"
	"regexp"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type Defuse struct {
	*elo.BaseEvent
	Subject     *elo.Player
	subjectTeam string
}

//"Jagger<19><STEAM_1:0:681607><CT>" triggered "Defused_The_Bomb"
var defusedrex = regexp.MustCompile(`"(.+)<(.+)><(.+)><(.+)>" triggered "Defused_The_Bomb"`)

func NewDefuseEvent(b *elo.BaseEvent) (e *Defuse) {
	if sm := defusedrex.FindStringSubmatch(b.Message); sm != nil {
		e = &Defuse{Subject: elo.GetPlayer(sm[1], sm[3]), subjectTeam: sm[4],
			BaseEvent: b}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Defuse) String() string {
	return fmt.Sprintf("Bomb defused by %s", e.Subject)
}
