package events

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
)

type MatchEnd struct {
	BaseEvent
	Duration time.Duration
	CtScore  int
	TScore   int
}

//Game Over: competitive default cs_office score 8:4 after 12 min
var matchendtrex = regexp.MustCompile(`^Game Over.+score (.+):(.+) after (.+) (.+)$`)

func NewMatchEndEvent(server *elo.Server, t time.Time, message string) (e *MatchEnd) {
	if sm := matchendtrex.FindStringSubmatch(message); sm != nil {
		cts, err1 := strconv.Atoi(sm[1])
		ts, err2 := strconv.Atoi(sm[2])
		d, err3 := time.ParseDuration(sm[3] + sm[4][0:1])
		if err1 != nil || err2 != nil || err3 != nil {
			log.Error("Malformed match end event: %s", message)
		} else {
			e = &MatchEnd{Duration: d, CtScore: cts, TScore: ts, BaseEvent: BaseEvent{Time: t, Server: server, Message: message}}
			log.Info("Created event: %+v", e)
		}
	}
	return e
}

func (e *MatchEnd) String() string {
	return fmt.Sprintf("Match end at %s, %d:%d, %s", e.Time.Format(time.RFC822Z), e.CtScore, e.TScore, e.Duration)
}
