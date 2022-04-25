package elo

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wlbr/commons/log"
)

type Emitter interface {
	WaitForProcessors()
	Loop()
	AddProcessor(p Processor)
	GetProcessor() []Processor
	AddFilter(f Filter)
	GetFilters() []Filter
}

func filter(em Emitter, m string) bool {
	for _, f := range em.GetFilters() {
		if f.Test(m) {
			return true
		}
	}
	return false
}

//================================

func shortenMessage(str string) (timestamp time.Time, message string, err error) {
	//str := string(buf)
	start := strings.IndexByte(str, ' ')
	if start == -1 {
		return timestamp, message, fmt.Errorf("Logline not in standard format. Logline: '%s'", str)
	}
	start++ //skip the initial space
	dend := start + 21
	if len(str) < dend || str[dend] != ':' {
		return timestamp, message, fmt.Errorf("Logline not in standard format, did not find end of date. Logline: '%s'", str)
	}
	if err == nil {
		layout := "01/02/2006 - 15:04:05"
		timestamp, err = time.Parse(layout, str[start:dend])
		if err != nil {
			log.Error("Could not parse event time, using <now>: %s", str[start:dend])
			timestamp = time.Now()
			err = nil
		}
		message = str[dend+2:]

	}
	return timestamp, message, err
}

//================================

type fileEmitter struct {
	wg            *sync.WaitGroup
	config        *Config
	procs         []Processor
	recordingfile *os.File
	wbuf          *bufio.Writer
	filters       []Filter
}

func NewFileEmitter(cfg *Config) *fileEmitter {
	e := &fileEmitter{config: cfg}
	e.wg = new(sync.WaitGroup)
	if cfg.Elo.RecorderFileName != "" {
		//f, err := os.Create(cfg.Elo.RecorderFileName)
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

func (em *fileEmitter) GetProcessor() []Processor {
	return em.procs
}

func (em *fileEmitter) AddProcessor(p Processor) {
	em.procs = append(em.procs, p)
	p.AddWaitGroup(em.wg)
	go p.Loop()
}

func (em *fileEmitter) AddFilter(f Filter) {
	em.filters = append(em.filters, f)
}

func (em *fileEmitter) GetFilters() []Filter {
	return em.filters
}

func (em *fileEmitter) Loop() {
	log.Debug("Starting file emitter loop.")
	f, err := os.Open(em.config.Elo.ImportFileName)
	if err != nil {
		log.Fatal("Error opening import CsLogFile '%s':  %s", em.config.Elo.ImportFileName, err)
	}

	scanner := bufio.NewScanner(f)
	lineno := 0
	server := &Server{IP: "fromFile"}
	for scanner.Scan() {
		lineno++
		buf := scanner.Text()
		t, m, err := shortenMessage(buf)
		if err != nil {
			log.Info("Ignoring line  %d. %v", lineno, err)
		} else {
			if em.config.Elo.RecorderFileName != "" {
				em.wbuf.WriteString(buf + "\n")
				em.wbuf.Flush()
			}
			if !filter(em, m) {
				for _, p := range em.procs {
					p.AddJob(NewBaseEvent(server, t, m))
				}
			}
		}
	}

	for _, p := range em.procs {
		p.AddJob(NewBaseEvent(server, time.Now(), "cselo:StopProcessing."))
	}

	//wait for the processors to stop
	em.WaitForProcessors()

	if err := scanner.Err(); err != nil {
		log.Fatal("Error scanning import CsLogFile '%s':  %s", em.config.Elo.ImportFileName, err)
	}
}

//================================

type udpEmitter struct {
	wg            *sync.WaitGroup
	config        *Config
	recordingfile *os.File
	procs         []Processor
	wbuf          *bufio.Writer
	filters       []Filter
}

func NewUdpEmitter(cfg *Config) *udpEmitter {
	e := new(udpEmitter)
	e.wg = new(sync.WaitGroup)
	e.config = cfg
	if cfg.Elo.RecorderFileName != "" {
		f, err := os.Create(cfg.Elo.RecorderFileName)
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

func (em *udpEmitter) AddProcessor(p Processor) {
	em.procs = append(em.procs, p)
	p.AddWaitGroup(em.wg)
	go p.Loop()
}

func (em *udpEmitter) GetProcessor() []Processor {
	return em.procs
}

func (em *udpEmitter) AddFilter(f Filter) {
	em.filters = append(em.filters, f)
}

func (em *udpEmitter) GetFilters() []Filter {
	return em.filters
}

func (em *udpEmitter) stopWorkers(server *Server) {
	for _, p := range em.procs {
		p.AddJob(NewBaseEvent(server, time.Now(), "cselo:StopProcessing."))
	}

	//wait for the processors to stop
	em.WaitForProcessors()
}

func (em *udpEmitter) Loop() {
	const protocol = "udp"
	port := em.config.Elo.Port
	server := &Server{IP: "fromUDP"}

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
			t, m, err := shortenMessage(sbuf)
			if err != nil {
				log.Warn("Ignoring line. %v", err)
				break
			} else {
				if em.config.Elo.RecorderFileName != "" {
					em.wbuf.WriteString(sbuf)
					em.wbuf.Flush()
				}
				if !filter(em, m) {
					for _, p := range em.procs {
						p.AddJob(NewBaseEvent(server, t, m))
					}
				}
			}

		}
	}

}

//================================
