package sinks

import (
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/events"
)

type inMemoryCounterSink struct {
	player           string
	allkills         []*events.Kill
	playerskills     []*events.Kill
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

func NewInMemoryCounterSink(cfg *elo.Config, playername string) (*inMemoryCounterSink, error) {
	return &inMemoryCounterSink{player: playername}, nil
}

func (s *inMemoryCounterSink) HandleKillEvent(e *events.Kill) {
	s.allkills = append(s.allkills, e)
	if s.player == e.Subject.Name {
		s.playerskills = append(s.playerskills, e)
	}
}

func (s *inMemoryCounterSink) HandleAssistEvent(e *events.Assist) {
	s.allassists = append(s.allassists, e)
	if s.player == e.Subject.Name {
		s.playersassists = append(s.playersassists, e)
	}
}
func (s *inMemoryCounterSink) HandleBlindedEvent(e *events.Blinded) {}
func (s *inMemoryCounterSink) HandleGrenadeEvent(e *events.Grenade) {}
func (s *inMemoryCounterSink) HandlePlantedEvent(e *events.Planted) {}

func (s *inMemoryCounterSink) HandleDefuseEvent(e *events.Defuse) {
	s.alldefuses = append(s.alldefuses, e)
	if s.player == e.Subject.Name {
		s.playersdefuses = append(s.playersdefuses, e)
	}
}

func (s *inMemoryCounterSink) HandleBombedEvent(e *events.Bombed) {
	s.allbombings = append(s.allbombings, e)
	if s.player == e.Subject.Name {
		s.playersbombings = append(s.playersbombings, e)
	}
}

func (s *inMemoryCounterSink) HandleHostageRescuedEvent(e *events.HostageRescued) {
	s.allrescues = append(s.allrescues, e)
	if s.player == e.Subject.Name {
		s.playersrescues = append(s.playersrescues, e)
	}
}

func (s *inMemoryCounterSink) HandleRoundStartEvent(e *events.RoundStart) {}
func (s *inMemoryCounterSink) HandleRoundEndEvent(e *events.RoundEnd)     {}
func (s *inMemoryCounterSink) HandleMatchStartEvent(e *events.MatchStart) {}
func (s *inMemoryCounterSink) HandleMatchEndEvent(e *events.MatchEnd)     {}
func (s *inMemoryCounterSink) HandleAccoladeEvent(e *events.Accolade)     {}
