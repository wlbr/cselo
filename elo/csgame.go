package elo

import (
	"sync"
	"time"
)

// type Round struct {
// 	ID       int64
// 	Match    *Match
// 	Duration time.Duration
// 	Start    time.Time
// 	End      time.Time
// }

// func (r *Round) CalcDuration() {
// 	if !r.End.IsZero() && !r.Start.IsZero() {
// 		r.Duration = r.End.Sub(r.Start)
// 	}
// }

type Match struct {
	ID              int64
	Server          *Server
	GameMode        string
	MapGroup        string
	MapFullName     string
	MapName         string
	ScoreA          int
	ScoreB          int
	Rounds          int
	Start           time.Time
	End             time.Time
	Duration        time.Duration
	playersNameLock sync.RWMutex
	playersByName   map[string]*Player
	playersIdLock   sync.RWMutex
	playersById     map[string]*Player
}

func (m *Match) AddPlayer(p *Player) {
	m.playersNameLock.Lock()
	if m.playersByName[p.Name] == nil {
		m.playersByName[p.Name] = p
	}
	m.playersNameLock.Unlock()
	m.playersIdLock.Lock()
	if m.playersById[p.SteamID] == nil {
		m.playersById[p.SteamID] = p
	}
	m.playersIdLock.Unlock()
}

func (m *Match) GetPlayerByName(name string) (p *Player) {
	m.playersNameLock.RLock()
	p = m.playersByName[name]
	m.playersNameLock.RUnlock()
	return p
}

func (m *Match) GetPlayerById(steamid string) (p *Player) {
	m.playersIdLock.RLock()
	p = m.playersById[steamid]
	m.playersIdLock.RUnlock()
	return p
}
