package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

//"Tina<217><BOT><CT>" assisted killing "Franzi<216><BOT><TERRORIST>"
//"Jagger<19><STEAM_1:0:681607><TERRORIST>" assisted killing "AHA<199><STEAM_1:1:689719><CT>"
var assistrex = regexp.MustCompile(`^"(.+)<(.+)><(.+)><(.+)>" assisted killing "(.+)<(.+)><(.+)><(.+)>"$`)

type Assist struct {
	BaseEvent
	Subject     *elo.Player
	subjectTeam string
	Object      *elo.Player
	objectTeam  string
}

func NewAssistEvent(server *elo.Server, t time.Time, message string) (e *Assist) {
	if sm := assistrex.FindStringSubmatch(message); sm != nil {
		e = &Assist{Subject: elo.GetPlayer(sm[1], sm[3]), subjectTeam: sm[4],
			Object: elo.GetPlayer(sm[5], sm[7]), objectTeam: sm[8],
			BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Assist) String() string {
	return fmt.Sprintf("Assist %s ==> %s", e.Subject, e.Object)
}
