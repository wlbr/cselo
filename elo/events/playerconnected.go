package events

import (
	"fmt"
	"regexp"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

//L 08/26/2021 - 18:00:55: "Dackel<21><STEAM_1:0:1770206><>" connected, address ""
//L 04/14/2022 - 18:27:16: "DorianHunter<39><STEAM_1:1:192746><>" connected, address ""
var connectrex = regexp.MustCompile(`"(.+)<(.+)><(.+)><(.*)>" connected, address "(.*)"`)

type PlayerConnected struct {
	*elo.BaseEvent
	Subject *elo.Player
	Address string
}

func NewPlayerConnectedEvent(b *elo.BaseEvent) (e *PlayerConnected) {
	if sm := connectrex.FindStringSubmatch(b.Message); sm != nil {
		e = &PlayerConnected{Subject: elo.GetPlayer(sm[1], sm[3]), Address: sm[5], BaseEvent: b}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *PlayerConnected) String() string {
	return fmt.Sprintf("Player %s connected from address %s", e.Subject, e.Address)
}
