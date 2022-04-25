package events

import (
	"fmt"
	"regexp"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type Bombed struct {
	*elo.BaseEvent
	Subject     *elo.Player
	subjectTeam string
}

func NewBombedEvent(b *elo.BaseEvent) (e *Bombed) {
	if sm := bombeddrex.FindStringSubmatch(b.Message); sm != nil {
		e = &Bombed{Subject: b.Server.LastPlanter, subjectTeam: sm[1],
			BaseEvent: b}
		log.Info("Created event: %+v", e)
	}
	return e
}

//Team "TERRORIST" triggered "SFUI_Notice_Target_Bombed" (CT "2") (T "2")
var bombeddrex = regexp.MustCompile(`Team "(.+)" triggered "SFUI_Notice_Target_Bombed" \((.+) "(.+)"\) \((.+) "(.+)"\)`)

func (e *Bombed) String() string {
	return fmt.Sprintf("Bomb exploded, planted by %s", e.Subject)
}
