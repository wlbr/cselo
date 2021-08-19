package sinks

import (
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/events"
)

type InMemoryCounterSink struct {
	Player           string
	AllKills         []*events.Kill
	PlayersKills     []*events.Kill
	allassists       []*events.Assist
	playersassists   []*events.Assist
	allplantings     []*events.Planted
	playersplantings []*events.Planted
	allbombings      []*events.Bombed
	playersbombings  []*events.Bombed
	alldefuses       []*events.Defuse
	playersdefuses   []*events.Defuse
	allrescues       []*events.HostageRescued
	playersrescues   []*events.HostageRescued
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
	s.allassists = append(s.allassists, e)
	if s.Player == e.Subject.Name {
		s.playersassists = append(s.playersassists, e)
	}
}
func (s *InMemoryCounterSink) HandleBlindedEvent(e *events.Blinded) {}
func (s *InMemoryCounterSink) HandleGrenadeEvent(e *events.Grenade) {}
func (s *InMemoryCounterSink) HandlePlantedEvent(e *events.Planted) {}

func (s *InMemoryCounterSink) HandleDefuseEvent(e *events.Defuse) {
	s.alldefuses = append(s.alldefuses, e)
	if s.Player == e.Subject.Name {
		s.playersdefuses = append(s.playersdefuses, e)
	}
}

func (s *InMemoryCounterSink) HandleBombedEvent(e *events.Bombed) {
	s.allbombings = append(s.allbombings, e)
	if s.Player == e.Subject.Name {
		s.playersbombings = append(s.playersbombings, e)
	}
}

func (s *InMemoryCounterSink) HandleHostageRescuedEvent(e *events.HostageRescued) {
	s.allrescues = append(s.allrescues, e)
	if s.Player == e.Subject.Name {
		s.playersrescues = append(s.playersrescues, e)
	}
}

func (s *InMemoryCounterSink) HandleRoundStartEvent(e *events.RoundStart) {}
func (s *InMemoryCounterSink) HandleRoundEndEvent(e *events.RoundEnd)     {}
func (s *InMemoryCounterSink) HandleMatchStartEvent(e *events.MatchStart) {}
func (s *InMemoryCounterSink) HandleMatchEndEvent(e *events.MatchEnd)     {}
func (s *InMemoryCounterSink) HandleAccoladeEvent(e *events.Accolade)     {}
