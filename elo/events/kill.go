package events

import (
	"fmt"
	"regexp"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type Kill struct {
	BaseEvent
	Subject         *elo.Player
	subjectPosition string
	subjectTeam     string
	Object          *elo.Player
	objectPosition  string
	objectTeam      string
	Weapon          string
	Headshot        bool
}

//"Madlen<221><BOT><CT>" [2650 -3117 -130] killed "Franzi<216><BOT><TERRORIST>" [3649 -3151 -48] with "hkp2000"
//"KiF Charlies Silence<16><STEAM_1:0:710013><TERRORIST>" [3878 -2315 -102] killed "Steffi<219><BOT><CT>" [2762 -4031 -142] with "sg556" (headshot)
var killrex = regexp.MustCompile(`"(.+)<(.+)><(.+)><(.+)>" \[(.+)\] killed "(.+)<(.+)><(.+)><(.+)>" \[(.+)\] with "(.+)"( \((headshot)\))?`)

func NewKillEvent(server *elo.Server, t time.Time, message string) (e *Kill) {
	if sm := killrex.FindStringSubmatch(message); sm != nil {
		headshot := len(sm) >= 14 && sm[13] == "headshot"

		e = &Kill{Subject: elo.GetPlayer(sm[1], sm[3]), subjectTeam: sm[4], subjectPosition: sm[5],
			Object: elo.GetPlayer(sm[6], sm[8]), objectTeam: sm[9], objectPosition: sm[10],
			Weapon: sm[11], Headshot: headshot,
			BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *Kill) String() (s string) {
	s = fmt.Sprintf("Kill %s ==> %s ", e.Subject, e.Object)
	if e.Headshot {
		s = s + "(headshot) "
	}
	s = s + e.Time.String()
	return s
}
