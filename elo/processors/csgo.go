package processors

import (
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
	"github.com/wlbr/cs-elo/elo/events"
	"github.com/wlbr/cs-elo/elo/sinks"
)

type CsgoLog struct {
	config  *elo.Config
	sinks   []sinks.Sink
	servers map[string]*elo.Server
}

func NewCsgoLogProcessor(cfg *elo.Config) *CsgoLog {
	p := &CsgoLog{config: cfg}
	p.servers = make(map[string]*elo.Server)

	return p
}

func (p *CsgoLog) AddSink(s sinks.Sink) {
	log.Info("Adding sink to CsgoLog processor: %#v", s)
	p.sinks = append(p.sinks, s)
}

func (p *CsgoLog) Dispatch(em elo.Emitter, srv *elo.Server, t time.Time, m string) {
	// srv, ok := p.servers[server]
	// if ok {
	// 	srv = &elo.Server{IP: server}
	// 	p.servers[server] = srv
	// }
	if e := events.NewKillEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleKillEvent(e)
		}
		return
	}
	if e := events.NewAssistEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleAssistEvent(e)
		}
		return
	}
	if e := events.NewBlindedEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleBlindedEvent(e)
		}
		return
	}
	if e := events.NewGrenadeEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleGrenadeEvent(e)
		}
		return
	}
	if e := events.NewBombedEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleBombedEvent(e)
		}
		return
	}
	if e := events.NewDefuseEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleDefuseEvent(e)
		}
		return
	}
	if e := events.NewHostageRescuedEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleHostageRescuedEvent(e)
		}
		return
	}
	if e := events.NewPlantedEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandlePlantedEvent(e)
		}
		return
	}
	if e := events.NewRoundStartEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleRoundStartEvent(e)
		}
		return
	}
	if e := events.NewRoundEndEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleRoundEndEvent(e)
		}
		return
	}
	if e := events.NewMatchStartEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleMatchStartEvent(e)
		}
		return
	}
	if e := events.NewMatchEndEvent(srv, t, m); e != nil {
		for _, s := range p.sinks {
			s.HandleMatchEndEvent(e)
		}
		return
	}

}
