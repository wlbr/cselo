package emitter

import (
	"bufio"
	"net"
	"os"
	"sync"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type udpEmitter struct {
	wg            *sync.WaitGroup
	config        *elo.Config
	recordingfile *os.File
	procs         []elo.Processor
	wbuf          *bufio.Writer
	filters       []elo.Filter
	recorder      *elo.Recorder
}

func NewUdpEmitter(cfg *elo.Config) *udpEmitter {
	e := new(udpEmitter)
	e.wg = new(sync.WaitGroup)
	e.config = cfg
	if cfg.Elo.RecorderFileName != "" {
		e.recorder = elo.NewRecorder(cfg, e.wg)
		go e.recorder.Loop()
	}
	if cfg.Elo.RecorderFileName != "" {
		f, err := os.OpenFile(cfg.Elo.RecorderFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		e.recordingfile = f
		if err != nil {
			log.Fatal("Could not create recorder file '%s': %s", cfg.Elo.RecorderFileName, err)
			cfg.FatalExit()
		} else {
			cfg.AddCleanUpFn(e.recordingfile.Close)
			e.wbuf = bufio.NewWriter(e.recordingfile)
		}
	}
	return e
}

func (em *udpEmitter) WaitForProcessors() {
	em.wg.Wait()
}

func (em *udpEmitter) AddProcessor(p elo.Processor) {
	em.procs = append(em.procs, p)
	//p.AddWaitGroup(em.wg)
	//go p.Loop()
}

func (em *udpEmitter) GetProcessor() []elo.Processor {
	return em.procs
}

func (em *udpEmitter) AddFilter(f elo.Filter) {
	em.filters = append(em.filters, f)
}

func (em *udpEmitter) GetFilters() []elo.Filter {
	return em.filters
}

func (em *udpEmitter) stopWorkers(server *elo.Server) {
	for _, p := range em.procs {
		p.AddJob(elo.NewBaseEvent(server, time.Now(), "cselo:StopProcessing."))
	}

	//wait for the processors to stop
	em.WaitForProcessors()
}

func (em *udpEmitter) Loop() {
	const protocol = "udp"
	port := em.config.Elo.Port
	server := elo.NewServer("fromUDP")

	//Build the address
	udpAddr, err := net.ResolveUDPAddr(protocol, ":"+port)
	if err != nil {
		log.Error("Error building address: %s", err)
		return
	}
	//Create the connection
	pc, err := net.ListenUDP(protocol, udpAddr)
	if err != nil {
		log.Error("Error opening connection: %s", err)
	}
	log.Info("Starting to listen on port %s", port)
	// pc, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		log.Fatal("%v", err)
	}
	defer pc.Close()

	defer em.stopWorkers(server)
	//the event loop
	for {
		buf := make([]byte, 1024)
		//n, addr, err := pc.ReadFromUDP(buf)
		n, _, err := pc.ReadFromUDP(buf)
		if err != nil {
			log.Error("Error during receiving: %v", err)
			continue
		} else {
			sbuf := string(buf[5 : n-1])
			t, m, err := elo.ShortenMessage(sbuf)
			if err != nil {
				log.Warn("Ignoring line. %v", err)
				break
			} else {
				if em.config.Elo.RecorderFileName != "" {
					end := ""
					if len(sbuf) > 0 && sbuf[len(sbuf)-1] != '\n' {
						end = "\n"
					}
					em.recorder.Record(sbuf + end)
					em.wbuf.WriteString(sbuf + end)

					em.wbuf.Flush()
				}
				if !elo.CheckFilter(em, m) {
					for _, p := range em.procs {
						p.AddJob(elo.NewBaseEvent(server, t, m))
					}
				}
			}

		}
	}

}
