package elo

import (
	"fmt"
	"time"
)

type Server struct {
	IP           string
	CurrentMatch *Match
	LastPlanter  *Player
	// CurrentRound *Round
}

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
	ID          int64
	Server      *Server
	GameMode    string
	MapGroup    string
	MapFullName string
	MapName     string
	ScoreA      int64
	ScoreB      int64
	Start       time.Time
	End         time.Time
	Duration    time.Duration
}

type Player struct {
	ID      int64
	Name    string
	SteamID string
}

type PlayersCache map[int]*Player

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

type Intervall struct {
	Start time.Time
	End   time.Time
}

func (i *Intervall) String() string {
	tfrmt := time.ANSIC
	return fmt.Sprintf("[%s - %s]", i.Start.Format(tfrmt), i.End.Format(tfrmt))
}

const Day = 24 * time.Hour
const Week = 7 * Day

func NewIntervall(start, end time.Time) *Intervall {
	return &Intervall{Start: start, End: end}
}

func IntervallLastXDays(x int) *Intervall {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 99, 999, time.UTC)
	closetostart := now.Add(Day * time.Duration(-x))
	start := time.Date(closetostart.Year(), closetostart.Month(), closetostart.Day(), 0, 0, 0, 0, time.UTC)
	return NewIntervall(start, end)
}

func IntervallLastWeek() *Intervall {
	return IntervallLastXDays(7)
}

func IntervallLastXYears(x int) *Intervall {
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 99, 999, time.UTC)
	start := time.Date(now.Year()-x, now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return NewIntervall(start, end)
}
