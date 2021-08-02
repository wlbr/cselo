package elo

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wlbr/commons/log"
)

type Player struct {
	ID        int64
	Name      string
	SteamID   string
	ProfileID string
}

type PlayersCache map[int]*Player

var players = make(map[string]*Player)

func GetPlayer(name, steamid string) (p *Player) {
	if p = players[steamid]; p == nil {
		p = &Player{Name: name, SteamID: steamid, ProfileID: SteamIdToProfileId(steamid)}
		players[steamid] = p
	}
	return p
}

func (p *Player) String() string {
	return p.Name + "-" + p.SteamID + "-" + p.ProfileID
}

//STEAM_1:0:681607
//STEAM_1:1:2102196
func SteamIdToProfileId(steamid string) (profileid string) {
	segments := strings.Split(steamid, ":")
	sid, err := strconv.Atoi(segments[2])
	if err != nil {
		log.Error("Cannot convert steamid '%s' error: %v", steamid, err)
		return ""
	}
	lid, err := strconv.Atoi(segments[1])
	if err != nil {
		log.Error("Cannot convert steamid '%s' error: %v", steamid, err)
		return ""
	}
	profileid = fmt.Sprintf("%d", 76561197960265728+(2*sid)+lid)
	return profileid
}
