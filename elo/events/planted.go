package events

import (
	"fmt"
	"regexp"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type Planted struct {
	*elo.BaseEvent
	Subject     *elo.Player
	subjectTeam string
}

//"AHA<18><STEAM_1:1:689719><TERRORIST>" triggered "Planted_The_Bomb"
var plantedrex = regexp.MustCompile(`"(.+)<(.+)><(.+)><(.+)>" triggered "Planted_The_Bomb"`)

func NewPlantedEvent(b *elo.BaseEvent) (e *Planted) {
	if sm := plantedrex.FindStringSubmatch(b.Message); sm != nil {
		pl := elo.GetPlayer(sm[1], sm[3])
		b.Server.LastPlanter = pl
		e = &Planted{Subject: pl, subjectTeam: sm[4],
			BaseEvent: b}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Planted) String() string {
	return fmt.Sprintf("Bomb planted by %s", e.Subject)
}
