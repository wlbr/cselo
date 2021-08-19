package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type Planted struct {
	BaseEvent
	Subject     *elo.Player
	subjectTeam string
}

//"AHA<18><STEAM_1:1:689719><TERRORIST>" triggered "Planted_The_Bomb"
var plantedrex = regexp.MustCompile(`^"(.+)<(.+)><(.+)><(.+)>" triggered "Planted_The_Bomb"$`)

func NewPlantedEvent(server *elo.Server, t time.Time, message string) (e *Planted) {
	if sm := plantedrex.FindStringSubmatch(message); sm != nil {
		pl := elo.GetPlayer(sm[1], sm[3])
		server.LastPlanter = pl
		e = &Planted{Subject: pl, subjectTeam: sm[4],
			BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Planted) String() string {
	return fmt.Sprintf("Bomb planted by %s", e.Subject)
}
