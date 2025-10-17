package elo

import (
	"fmt"
	"sync"
	"time"

	"github.com/wlbr/commons/log"
)

type Server struct {
	m            sync.RWMutex
	IP           string
	currentMatch *Match
	lastMatch    *Match
	LastPlanter  *Player
	// CurrentRound *Round
}

func NewServer(ip string) *Server {
	log.Debug("NewServer created: ip=%s", ip)
	return &Server{IP: ip}
}

func (s *Server) SetCurrentMatch(m *Match) {
	log.Info("Setting current match for server %s: %+v", s.IP, m)
	s.m.Lock()
	defer s.m.Unlock()
	s.lastMatch = s.currentMatch
	s.currentMatch = m
}

func (s *Server) CurrentMatch() *Match {
	s.m.RLock()
	defer s.m.RUnlock()

	if s.currentMatch == nil {
		log.Debug("CurrentMatch is nil for server %s", s.IP)
		return nil
	}
	return s.currentMatch
}

func (s *Server) LastMatch() *Match {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.lastMatch
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
