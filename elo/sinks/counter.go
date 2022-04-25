package sinks

import (
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/events"
)

type InMemoryCounterSink struct {
	Player           string
	AllKills         []*events.Kill
	PlayersKills     []*events.Kill
	Allassists       []*events.Assist
	Playersassists   []*events.Assist
	Allplantings     []*events.Planted
	Playersplantings []*events.Planted
	Allbombings      []*events.Bombed
	Playersbombings  []*events.Bombed
	Alldefuses       []*events.Defuse
	Playersdefuses   []*events.Defuse
	Allrescues       []*events.HostageRescued
	Playersrescues   []*events.HostageRescued
	Matches          []*elo.Match
}

func NewInMemoryCounterSink(cfg *elo.Config, playername string) (*InMemoryCounterSink, error) {
	return &InMemoryCounterSink{Player: playername}, nil
}

func (s *InMemoryCounterSink) HandleKillEvent(e *events.Kill) {
	s.AllKills = append(s.AllKills, e)
	if s.Player == e.Subject.Name {
		s.PlayersKills = append(s.PlayersKills, e)
	}
}

func (s *InMemoryCounterSink) HandleAssistEvent(e *events.Assist) {
	s.Allassists = append(s.Allassists, e)
	if s.Player == e.Subject.Name {
		s.Playersassists = append(s.Playersassists, e)
	}
}
func (s *InMemoryCounterSink) HandleBlindedEvent(e *events.Blinded) {}
func (s *InMemoryCounterSink) HandleGrenadeEvent(e *events.Grenade) {}

func (s *InMemoryCounterSink) HandlePlantedEvent(e *events.Planted) {
	s.Allplantings = append(s.Allplantings, e)
	if s.Player == e.Subject.Name {
		s.Playersplantings = append(s.Playersplantings, e)
	}
}

func (s *InMemoryCounterSink) HandleDefuseEvent(e *events.Defuse) {
	s.Alldefuses = append(s.Alldefuses, e)
	if s.Player == e.Subject.Name {
		s.Playersdefuses = append(s.Playersdefuses, e)
	}
}

func (s *InMemoryCounterSink) HandleBombedEvent(e *events.Bombed) {
	s.Allbombings = append(s.Allbombings, e)
	if s.Player == e.Subject.Name {
		s.Playersbombings = append(s.Playersbombings, e)
	}
}

func (s *InMemoryCounterSink) HandleHostageRescuedEvent(e *events.HostageRescued) {
	s.Allrescues = append(s.Allrescues, e)
	if s.Player == e.Subject.Name {
		s.Playersrescues = append(s.Playersrescues, e)
	}
}

func (s *InMemoryCounterSink) HandleMatchEndEvent(e *events.MatchEnd) {
	m := e.Server.CurrentMatch
	s.Matches = append(s.Matches, m)
}

func (s *InMemoryCounterSink) HandleMatchStatusEvent(e *events.MatchStatus)             {}
func (s *InMemoryCounterSink) HandleMatchStartEvent(e *events.MatchStart)               {}
func (s *InMemoryCounterSink) HandleServerHibernationEvent(e *events.ServerHibernation) {}
func (s *InMemoryCounterSink) HandleRoundStartEvent(e *events.RoundStart)               {}
func (s *InMemoryCounterSink) HandleRoundEndEvent(e *events.RoundEnd)                   {}
func (s *InMemoryCounterSink) HandleAccoladeEvent(e *events.Accolade)                   {}
func (s *InMemoryCounterSink) HandleMatchCleanUpEvent(e *events.MatchCleanUp)           {}
func (s *InMemoryCounterSink) HandlePlayerConnectedEvent(e *events.PlayerConnected)     {}
