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

func (p *Player) String() string {
	return p.Name + "-" + p.SteamID + "-" + p.ProfileID
}

type PlayersCache map[int]*Player

var players = make(map[string]*Player)

type PlayerLookupError struct {
	description string
	name        string
	candidates  []*Player
}

func (e *PlayerLookupError) Error() string {
	return fmt.Sprintf("%s '%s': %v", e.description, e.name, e.candidates)
}

func GetPlayer(name, steamid string) (p *Player) {
	if p = players[steamid]; p == nil {
		p = &Player{Name: name, SteamID: steamid, ProfileID: SteamIdToProfileId(steamid)}
		players[steamid] = p
	}
	return p
}

func GetPlayerByName(name string) (p *Player, e *PlayerLookupError) {
	var cands []*Player
	for _, pl := range players {
		if pl.Name == name {
			cands = append(cands, pl)
		}
	}

	switch len(cands) {
	case 0:
		e = &PlayerLookupError{description: "Did not finy any player with name", name: name, candidates: cands}
	case 1:
		p = cands[0]
	default:
		e = &PlayerLookupError{description: "Found more than one player with name", name: name, candidates: cands}
	}
	return p, e
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
