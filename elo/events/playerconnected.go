package events

import (
	"fmt"
	"regexp"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

// CS2 http:
//
//		DorianHunter<1><[U:1:385493]><> connected, address 87.122.5.65:60704
//	 DorianHunter<1><[U:1:385493]><> connected, address 87.122.5.65:60704
//
// 11/04/2023 - 15:41:55.798 - "Jagger<0><[U:1:1363214]><>" connected, address "172.17.0.1:45612"
// CS2 logfile:
// L 10/26/2023 - 11:59:04: "Jagger<0><[U:1:1363214]><>" connected, address "172.17.0.1:50390"
// CSGO:
// L 08/26/2021 - 18:00:55: "Dackel<21><STEAM_1:0:1770206><>" connected, address ""
// L 04/14/2022 - 18:27:16: "DorianHunter<39><STEAM_1:1:192746><>" connected, address ""
var connectrex = regexp.MustCompile(`"?(.+)<(.+)><(.+)><(.*)>"? connected, address "?(.*)"?`)

type PlayerConnected struct {
	*elo.BaseEvent
	Subject *elo.Player
	Address string
}

func NewPlayerConnectedEvent(b *elo.BaseEvent) (e *PlayerConnected) {
	if sm := connectrex.FindStringSubmatch(b.Message); sm != nil {
		player := elo.GetPlayer(sm[1], sm[3])
		if player == nil {
			log.Error("NewPlayerConnectedEvent:Cannot determine player for event message '%s': '%s'", b.Time, b.Message)
			return nil
		}
		e = &PlayerConnected{Subject: player, Address: sm[5], BaseEvent: b}
		log.Info("Created event: %+v", e)
	}
	return e
}

func (e *PlayerConnected) String() string {
	return fmt.Sprintf("Player %s connected from address %s", e.Subject, e.Address)
}
