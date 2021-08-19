package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type Bombed struct {
	BaseEvent
	Subject     *elo.Player
	subjectTeam string
}

func NewBombedEvent(server *elo.Server, t time.Time, message string) (e *Bombed) {
	if sm := bombeddrex.FindStringSubmatch(message); sm != nil {
		e = &Bombed{Subject: server.LastPlanter, subjectTeam: sm[1],
			BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %+v", e)
	}
	return e
}

//Team "TERRORIST" triggered "SFUI_Notice_Target_Bombed" (CT "2") (T "2")
var bombeddrex = regexp.MustCompile(`Team "(.+)" triggered "SFUI_Notice_Target_Bombed" \((.+) "(.+)"\) \((.+) "(.+)"\)`)

func (e *Bombed) String() string {
	return fmt.Sprintf("Bomb exploded, planted by %s", e.Subject)
}
