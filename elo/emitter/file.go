package emitter

import (
	"bufio"
	"os"
	"sync"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type fileEmitter struct {
	wg            *sync.WaitGroup
	config        *elo.Config
	procs         []elo.Processor
	recordingfile *os.File
	wbuf          *bufio.Writer
	filters       []elo.Filter
}

func NewFileEmitter(cfg *elo.Config) *fileEmitter {
	e := &fileEmitter{config: cfg}
	e.wg = new(sync.WaitGroup)
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

func (em *fileEmitter) WaitForProcessors() {
	em.wg.Wait()
}

func (em *fileEmitter) GetProcessor() []elo.Processor {
	return em.procs
}

func (em *fileEmitter) AddProcessor(p elo.Processor) {
	em.procs = append(em.procs, p)
	p.AddWaitGroup(em.wg)
	//go p.Loop()
}

func (em *fileEmitter) AddFilter(f elo.Filter) {
	em.filters = append(em.filters, f)
}

func (em *fileEmitter) GetFilters() []elo.Filter {
	return em.filters
}

func (em *fileEmitter) Loop() {
	log.Info("Starting file emitter loop.")
	f, err := os.Open(em.config.Elo.ImportFileName)
	if err != nil {
		log.Fatal("Error opening import CsLogFile '%s':  %s", em.config.Elo.ImportFileName, err)
	}

	scanner := bufio.NewScanner(f)
	lineno := 0
	server := elo.NewServer("fromFile")
	for scanner.Scan() {
		lineno++
		buf := scanner.Text()
		t, m, err := elo.ShortenMessage(buf)
		if err != nil {
			log.Info("Ignoring line  %d. %v", lineno, err)
		} else {
			if em.config.Elo.RecorderFileName != "" {
				em.wbuf.WriteString(buf + "\n")
				em.wbuf.Flush()
			}
			if !elo.CheckFilter(em, m) {
				for _, p := range em.procs {
					p.AddJob(elo.NewBaseEvent(server, t, m))
				}
			}
		}
	}

	for _, p := range em.procs {
		p.AddJob(elo.NewBaseEvent(server, time.Now(), "cselo:StopProcessing."))
	}

	//wait for the processors to stop
	em.WaitForProcessors()

	if err := scanner.Err(); err != nil {
		log.Fatal("Error scanning import CsLogFile '%s':  %s", em.config.Elo.ImportFileName, err)
	}
}
