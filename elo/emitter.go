package elo

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
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

// CS2 http:
// 11/04/2023 - 15:41:55.798 - "Jagger<0><[U:1:1363214]><>" connected, address "172.17.0.1:45612"
// CS2 logfile:
// L 10/26/2023 - 11:59:04: "Jagger<0><[U:1:1363214]><>" connected, address "172.17.0.1:50390"
// CSGO:
// L 08/26/2021 - 18:00:55: "Dackel<21><STEAM_1:0:1770206><>" connected, address ""
// L 04/14/2022 - 18:27:16: "DorianHunter<39><STEAM_1:1:192746><>" connected, address ""

// func shortenMessage(str string) (timestamp time.Time, message string, err error) {
// 	strings.TrimLeft(str, " L")
// 	start := 0
// 	dend := start + 25
// 	if len(str) < dend || str[dend] != ' ' {
// 		return timestamp, message, fmt.Errorf("Logline not in standard format, did not find end of date. Logline: '%s'", str)
// 	}
// 	if err == nil {
// 		layout := "01/02/2006 - 15:04:05.000"
// 		timestamp, err = time.Parse(layout, str[start:dend])
// 		if err != nil {
// 			log.Error("Could not parse event time, using <now>: %s", str[start:dend])
// 			timestamp = time.Now()
// 			err = nil
// 		}
// 		message = str[dend+3:]
// 	}

// 	return timestamp, message, err
// }

// var shortenrex = regexp.MustCompile(`(?Um)L |(.+) - (.+)(:| -) (.*)$`)
// var shortenrex = regexp.MustCompile(`(L )?(.+) - (\d\d:\d\d:(\d|\.)+)(:| -) (.*)$`)
var shortenrex = regexp.MustCompile(`(.+) - (\d\d:\d\d:(\d|\.)+)(:| -) (.*)$`)

func shortenMessage(str string) (timestamp time.Time, message string, err error) {

	str = strings.TrimRight(str, "\n")
	str = strings.TrimLeft(str, "L ")

	if sm := shortenrex.FindStringSubmatch(str); sm != nil {
		layout := "01/02/2006 - 15:04:05"
		if len(sm[2]) > 8 {
			layout += ".000"
		}
		timestamp, err = time.Parse(layout, sm[1]+" - "+sm[2])
		if err != nil {
			log.Error("Could not parse event time, using <now>: %s", sm[1]+" - "+sm[2])
			timestamp = time.Now()
			err = nil
		}
		message = sm[5]

		return timestamp, message, err
	} else {
		return timestamp, message, fmt.Errorf("Logline not in standard format. Logline: '%s'", str)
	}
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

type httpEmitter struct {
	wg            *sync.WaitGroup
	config        *Config
	recordingfile *os.File
	procs         []Processor
	wbuf          *bufio.Writer
	filters       []Filter
}

func NewHttpEmitter(cfg *Config) *httpEmitter {
	e := new(httpEmitter)
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

func (em *httpEmitter) WaitForProcessors() {
	em.wg.Wait()
}

func (em *httpEmitter) AddProcessor(p Processor) {
	em.procs = append(em.procs, p)
	p.AddWaitGroup(em.wg)
	go p.Loop()
}

func (em *httpEmitter) GetProcessor() []Processor {
	return em.procs
}

func (em *httpEmitter) AddFilter(f Filter) {
	em.filters = append(em.filters, f)
}

func (em *httpEmitter) GetFilters() []Filter {
	return em.filters
}

func (em *httpEmitter) stopWorkers(server *Server) {
	for _, p := range em.procs {
		p.AddJob(NewBaseEvent(server, time.Now(), "cselo:StopProcessing."))
	}

	//wait for the processors to stop
	em.WaitForProcessors()
}

func (em *httpEmitter) Loop() {
	const protocol = "http"
	port := "42820" //em.config.Elo.Port
	server := &Server{IP: "fromHttp"}
	defer em.stopWorkers(server)
	handler := &csLogHandler{emitter: em, server: server}

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    1 * time.Second,
		WriteTimeout:   1 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Info("Starting to listen on port %s", port)
	err := s.ListenAndServe()

	if err != nil {
		log.Error("Error opening http server: %s", err)
	}
}

type csLogHandler struct {
	emitter *httpEmitter
	server  *Server
}

func (h csLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	var err error

	if r.Body == nil {
		log.Info("Empty request body. Url: %s", r.URL)
	} else {
		buf, err = io.ReadAll(r.Body)
		//fmt.Println(string(buf))
		if err != nil {
			log.Error("Problem reading request body: %v", err)
		} else {
			sbuf := string(buf)
			t, m, err := shortenMessage(sbuf)
			if err != nil {
				log.Warn("Ignoring line. %v", err)
				return
			} else {
				if h.emitter.config.Elo.RecorderFileName != "" {
					h.emitter.wbuf.WriteString(sbuf)
					h.emitter.wbuf.Flush()
				}
				if !filter(h.emitter, m) {
					for _, p := range h.emitter.procs {
						p.AddJob(NewBaseEvent(h.server, t, m))
					}
				}
			}
		}
	}
}

//================================
