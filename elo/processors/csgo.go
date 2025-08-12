package processors

import (
	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/events"
	"github.com/wlbr/cselo/elo/sinks"
)

type CsgoLog struct {
	incoming chan *elo.BaseEvent
	config   *elo.Config
	sinks    []sinks.Sink
	servers  map[string]*elo.Server
}

func NewCsgoLogProcessor(cfg *elo.Config) *CsgoLog {
	p := &CsgoLog{config: cfg}
	//p.m = &sync.Mutex{}
	p.servers = make(map[string]*elo.Server)
	p.incoming = make(chan *elo.BaseEvent) //, cfg.Elo.BufferSize)
	return p
}

// func (p *CsgoLog) AddWaitGroup(wg *sync.WaitGroup) {
// 	p.wg = wg
// }

func (p *CsgoLog) AddSink(s sinks.Sink) {
	log.Info("Adding sink to CsgoLog processor: %#v", s)
	p.sinks = append(p.sinks, s)
}

func (p *CsgoLog) AddJob(b *elo.BaseEvent) {
	p.incoming <- b
}

func (p *CsgoLog) Loop() {
	// p.wg.Add(1)
	// defer p.wg.Done()
	for {
		e := <-p.incoming
		p.process(e)
		if e.Message == "cselo:StopProcessing." {
			break
		}
	}
	defer log.Info("Finishing processor")
}

// func (p *CsgoLog) Dispatch(em elo.Emitter, b.Server *elo.Server, t time.Time, m string) {
func (p *CsgoLog) process(b *elo.BaseEvent) {
	log.Info("Processing event: %s - %v", b.Message, b.Server.CurrentMatch())
	if b.Server.CurrentMatch() == nil {
		log.Info("No current match, creating a new one.")
		match := elo.NewMatch("unknown", "unknown", b.Time, b.Server)
		b.Server.SetCurrentMatch(match)
		mse := &events.MatchStart{
			BaseEvent: &elo.BaseEvent{
				Server:  b.Server,
				Time:    b.Time,
				Message: "Missed MatchStart, guessing new one.",
			},
			MapFullName: "unknown",
			MapName:     "unknown",
		}
		for _, sink := range p.sinks {
			sink.HandleMatchStartEvent(mse)
		}
	}
	if e := events.NewKillEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleKillEvent(e)
		}
		return
	}
	if e := events.NewAssistEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleAssistEvent(e)
		}
		return
	}
	if e := events.NewBlindedEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleBlindedEvent(e)
		}
		return
	}
	if e := events.NewGrenadeEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleGrenadeEvent(e)
		}
		return
	}
	if e := events.NewBombedEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleBombedEvent(e)
		}
		return
	}
	if e := events.NewDefuseEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleDefuseEvent(e)
		}
		return
	}
	if e := events.NewHostageRescuedEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleHostageRescuedEvent(e)
		}
		return
	}
	if e := events.NewPlantedEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandlePlantedEvent(e)
		}
		return
	}
	if e := events.NewRoundStartEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleRoundStartEvent(e)
		}
		return
	}
	if e := events.NewRoundEndEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleRoundEndEvent(e)
		}
		return
	}
	//oldmatch := b.Server.CurrentMatch
	if e := events.NewMatchStartEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleMatchStartEvent(e)
		}
		//c := events.NewMatchCleanUpEvent(b.Server, e.Time, fmt.Sprintf("MatchCleanUp: Check for empty match %d", oldmatch.ID), oldmatch)
		c := events.NewMatchCleanUpEvent(b, b.Server.LastMatch())
		for _, s := range p.sinks {
			s.HandleMatchCleanUpEvent(c)
		}
		return
	}
	if e := events.NewMatchEndEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleMatchEndEvent(e)
		}
		return
	}
	if e := events.NewServerHibernationEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleServerHibernationEvent(e)
		}
		c := events.NewMatchCleanUpEvent(b, b.Server.LastMatch())
		for _, s := range p.sinks {
			s.HandleMatchCleanUpEvent(c)
		}
		return
	}
	if e := events.NewMatchStatusEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleMatchStatusEvent(e)
		}
		return
	}
	if e := events.NewAccoladeEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandleAccoladeEvent(e)
		}
		return
	}
	if e := events.NewPlayerConnectedEvent(b); e != nil {
		for _, s := range p.sinks {
			s.HandlePlayerConnectedEvent(e)
		}
		return
	}

}
