package elo

import "time"

type Server struct {
	IP           string
	CurrentMatch *Match
	CurrentRound *Round
}

type Round struct {
	Start      time.Time
	End        time.Time
	Players    []Player
	WinnerTeam string
}

type Match struct {
	Start      time.Time
	End        time.Time
	Players    []Player
	WinnerTeam string
	CtScore    int
	TScore     int
}

type Player struct {
	ID      int64
	Name    string
	SteamID string
}

var players = make(map[string]*Player)

func GetPlayer(name, steamid string) (p *Player) {
	if p = players[steamid]; p == nil {
		p = &Player{Name: name, SteamID: steamid}
		players[steamid] = p
	}
	return p
}

func (p *Player) String() string {
	return p.Name
}
