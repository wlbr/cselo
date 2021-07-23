package sources

import (
	"bufio"
	"sync"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
	"github.com/wlbr/cs-elo/elo/events"
)

type printer struct {
	config *elo.Config
	w      *bufio.Writer
	mu     sync.Mutex
}

func Newprinter(cfg *elo.Config) (*printer, error) {
	return &printer{config: cfg, w: bufio.NewWriter(cfg.Elo.OutputFile)}, nil
}

func (s *printer) printToFile(e elo.Event) {
	log.Debug("Printing  %v", e)
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.w.WriteString(e.String() + "\n")
	if err != nil {
		log.Fatal("Error writing to file: %v  error: %v", e, err)
	}

	s.w.Flush()
}

func (s *printer) HandleKillEvent(e *events.Kill)                     { s.printToFile(e) }
func (s *printer) HandleAssistEvent(e *events.Assist)                 { s.printToFile(e) }
func (s *printer) HandleBlindedEvent(e *events.Blinded)               { s.printToFile(e) }
func (s *printer) HandleGrenadeEvent(e *events.Grenade)               { s.printToFile(e) }
func (s *printer) HandlePlantedEvent(e *events.Planted)               { s.printToFile(e) }
func (s *printer) HandleDefuseEvent(e *events.Defuse)                 { s.printToFile(e) }
func (s *printer) HandleBombedEvent(e *events.Bombed)                 { s.printToFile(e) }
func (s *printer) HandleHostageRescuedEvent(e *events.HostageRescued) { s.printToFile(e) }
func (s *printer) HandleRoundStartEvent(e *events.RoundStart)         { s.printToFile(e) }
func (s *printer) HandleRoundEndEvent(e *events.RoundEnd)             { s.printToFile(e) }
func (s *printer) HandleMatchStartEvent(e *events.MatchStart)         { s.printToFile(e) }
func (s *printer) HandleMatchEndEvent(e *events.MatchEnd)             { s.printToFile(e) }
