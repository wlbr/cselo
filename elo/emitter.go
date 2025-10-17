package elo

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/wlbr/commons/log"
)

type Emitter interface {
	WaitForProcessors()
	Loop()
	AddProcessor(p Processor)
	GetProcessor() []Processor
	AddFilter(f Filter)
	GetFilters() []Filter
}

//================================

// CS2 http:
// 11/04/2023 - 15:41:55.798 - "Jagger<0><[U:1:1363214]><>" connected, address "172.17.0.1:45612"
// CS2 logfile:
// L 10/26/2023 - 11:59:04: "Jagger<0><[U:1:1363214]><>" connected, address "172.17.0.1:50390"
// CSGO:
// L 08/26/2021 - 18:00:55: "Dackel<21><STEAM_1:0:1770206><>" connected, address ""
// L 04/14/2022 - 18:27:16: "DorianHunter<39><STEAM_1:1:192746><>" connected, address ""

var shortenrex = regexp.MustCompile(`"?(.+) - (\d\d:\d\d:(\d|\.)+)(:| -) (.*)"?$`)

func ShortenMessage(str string) (timestamp time.Time, message string, err error) {

	str = strings.TrimRight(str, "\n")
	str = strings.TrimLeft(str, "L ")

	if sm := shortenrex.FindStringSubmatch(str); sm != nil {
		layout := "01/02/2006 - 15:04:05"
		if len(sm[2]) > 8 {
			layout += ".000"
		}
		timestamp, err = time.Parse(layout, sm[1]+" - "+sm[2])
		if err != nil {
			log.Error("Could not parse event time, using <now>: %s", sm[1]+" - "+sm[2])
			timestamp = time.Now()
			err = nil
		}
		message = sm[5]

		return timestamp, message, err
	} else {
		return timestamp, message, fmt.Errorf("Logline not in standard format. Logline: '%s'", str)
	}
}
