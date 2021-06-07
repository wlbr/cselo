package sinks

import (
	"bufio"
	"sync"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
	"github.com/wlbr/cs-elo/elo/events"
)

type printerSink struct {
	config *elo.Config
	w      *bufio.Writer
	mu     sync.Mutex
}

func NewPrinterSink(cfg *elo.Config) (*printerSink, error) {
	return &printerSink{config: cfg, w: bufio.NewWriter(cfg.Elo.OutputFile)}, nil
}

func (s *printerSink) printToFile(e elo.Event) {
	log.Debug("Printing  %v", e)
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.w.WriteString(e.String() + "\n")
	if err != nil {
		log.Fatal("Error writing to file: %v  error: %v", e, err)
	}

	s.w.Flush()
}

func (s *printerSink) HandleKillEvent(e *events.Kill)                     { s.printToFile(e) }
func (s *printerSink) HandleAssistEvent(e *events.Assist)                 { s.printToFile(e) }
func (s *printerSink) HandleBlindedEvent(e *events.Blinded)               { s.printToFile(e) }
func (s *printerSink) HandleGrenadeEvent(e *events.Grenade)               { s.printToFile(e) }
func (s *printerSink) HandlePlantedEvent(e *events.Planted)               { s.printToFile(e) }
func (s *printerSink) HandleDefuseEvent(e *events.Defuse)                 { s.printToFile(e) }
func (s *printerSink) HandleBombedEvent(e *events.Bombed)                 { s.printToFile(e) }
func (s *printerSink) HandleHostageRescuedEvent(e *events.HostageRescued) { s.printToFile(e) }
func (s *printerSink) HandleRoundStartEvent(e *events.RoundStart)         { s.printToFile(e) }
func (s *printerSink) HandleRoundEndEvent(e *events.RoundEnd)             { s.printToFile(e) }
func (s *printerSink) HandleMatchStartEvent(e *events.MatchStart)         { s.printToFile(e) }
func (s *printerSink) HandleMatchEndEvent(e *events.MatchEnd)             { s.printToFile(e) }
