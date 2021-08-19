package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type HostageRescued struct {
	BaseEvent
	Subject     *elo.Player
	subjectTeam string
}

//"Jagger<19><STEAM_1:0:681607><CT>" triggered "Rescued_A_Hostage"
var rescuedrex = regexp.MustCompile(`^"(.+)<(.+)><(.+)><(.+)>" triggered "Rescued_A_Hostage"$`)

func NewHostageRescuedEvent(server *elo.Server, t time.Time, message string) (e *HostageRescued) {
	if sm := rescuedrex.FindStringSubmatch(message); sm != nil {
		e = &HostageRescued{Subject: elo.GetPlayer(sm[1], sm[3]), subjectTeam: sm[4],
			BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *HostageRescued) String() string {
	return fmt.Sprintf("Hostage rescued by %s", e.Subject)
}
