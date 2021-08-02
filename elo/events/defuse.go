package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

type Defuse struct {
	BaseEvent
	Subject     *elo.Player
	subjectTeam string
}

//"Jagger<19><STEAM_1:0:681607><CT>" triggered "Defused_The_Bomb"
var defusedrex = regexp.MustCompile(`^"(.+)<(.+)><(.+)><(.+)>" triggered "Defused_The_Bomb"$`)

func NewDefuseEvent(server *elo.Server, t time.Time, message string) (e *Defuse) {
	if sm := defusedrex.FindStringSubmatch(message); sm != nil {
		e = &Defuse{Subject: elo.GetPlayer(sm[1], sm[3]), subjectTeam: sm[4],
			BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Defuse) String() string {
	return fmt.Sprintf("Bomb defused by %s", e.Subject)
}
