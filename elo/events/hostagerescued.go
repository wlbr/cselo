package events

import (
	"fmt"
	"regexp"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type HostageRescued struct {
	*elo.BaseEvent
	Subject     *elo.Player
	subjectTeam string
}

// CSGO
// "Jagger<19><STEAM_1:0:681607><CT>" triggered "Rescued_A_Hostage"
// CS2
// "Jagger<0><[U:1:1363214]><CT>" triggered "Rescued_A_Hostage"
var rescuedrex = regexp.MustCompile(`"(.+)<(.+)><(.+)><(.+)>" triggered "Rescued_A_Hostage"`)

func NewHostageRescuedEvent(b *elo.BaseEvent) (e *HostageRescued) {
	if sm := rescuedrex.FindStringSubmatch(b.Message); sm != nil {
		e = &HostageRescued{Subject: elo.GetPlayer(sm[1], sm[3]), subjectTeam: sm[4],
			BaseEvent: b}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *HostageRescued) String() string {
	return fmt.Sprintf("Hostage rescued by %s", e.Subject)
}
