package elo

import (
	"sync"
	"time"

	"github.com/wlbr/commons/log"
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
	idLock          sync.RWMutex
	setterLock      sync.RWMutex
	mapFullNameLock sync.RWMutex
	playersIdLock   sync.RWMutex
	playersNameLock sync.RWMutex

	idx           int64
	Server        *Server
	GameMode      string
	MapGroup      string
	mapFullNamex  string
	MapName       string
	ScoreA        int
	ScoreB        int
	Rounds        int
	Start         time.Time
	End           time.Time
	Duration      time.Duration
	Completed     bool
	playersByName map[string]*Player
	playersById   map[string]*Player
}

func NewMatch(MapFullName string, MapName string, Start time.Time, Server *Server) *Match {
	m := &Match{
		Server:        Server,
		Start:         Start,
		mapFullNamex:  MapFullName,
		MapName:       MapName,
		playersByName: make(map[string]*Player),
		playersById:   make(map[string]*Player),
	}
	return m
}

func (m *Match) Set(GameMode string, MapGroup string, MapFullName string, MapName string, ScoreA int, ScoreB int, Completed bool, End time.Time, Duration time.Duration) {
	m.setterLock.Lock()
	defer m.setterLock.Unlock()
	m.GameMode = GameMode
	m.MapGroup = MapGroup
	m.mapFullNamex = MapFullName
	m.MapName = MapName
	m.ScoreA = ScoreA
	m.ScoreB = ScoreB
	m.Completed = Completed
	m.End = End
	m.Duration = Duration
}

func (m *Match) SetStatusAttributes(mapFullName, mapName string, scoreA, scoreB, rounds int) {
	log.Warn("In SetStatusAttributes L1: mutex=%v", m)
	m.setterLock.Lock()
	log.Warn("In SetStatusAttributes L3")
	defer m.setterLock.Unlock()
	m.mapFullNamex = mapFullName
	m.MapName = mapName
	m.ScoreA = scoreA
	m.ScoreB = scoreB
	m.Rounds = rounds
}

func (m *Match) ID() int64 {
	m.idLock.RLock()
	defer m.idLock.RUnlock()
	return m.idx
}

func (m *Match) SetID(id int64) {
	m.idLock.Lock()
	defer m.idLock.Unlock()
	m.idx = id
}

func (m *Match) MapFullName() string {
	m.mapFullNameLock.RLock()
	defer m.mapFullNameLock.RUnlock()
	return m.mapFullNamex
}

// func (m *Match) SetMapFullName(name string) {
// 	m.mapFullNameLock.Lock()
// 	defer m.mapFullNameLock.Unlock()
// 	m.mapFullName = name
// }

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
	defer m.playersNameLock.RUnlock()
	p = m.playersByName[name]
	return p
}

func (m *Match) GetPlayerById(steamid string) (p *Player) {
	m.playersIdLock.RLock()
	defer m.playersIdLock.RUnlock()
	p = m.playersById[steamid]
	return p
}
