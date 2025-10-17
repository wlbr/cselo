package emitter

import (
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
)

type httpEmitter struct {
	//m             *sync.Mutex
	wg            *sync.WaitGroup
	config        *elo.Config
	recordingfile *os.File
	procs         []elo.Processor
	filters       []elo.Filter
	recorder      *elo.Recorder
}

func NewHttpEmitter(cfg *elo.Config) *httpEmitter {
	e := new(httpEmitter)
	e.wg = new(sync.WaitGroup)
	if cfg.Elo.RecorderFileName != "" {
		e.recorder = elo.NewRecorder(cfg, e.wg)
		go e.recorder.Loop()
	}
	e.config = cfg
	if cfg.Elo.RecorderFileName != "" {
		f, err := os.OpenFile(cfg.Elo.RecorderFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		e.recordingfile = f
		if err != nil {
			log.Fatal("Could not create recorder file '%s': %s", cfg.Elo.RecorderFileName, err)
			cfg.FatalExit()
		} else {
			cfg.AddCleanUpFn(e.recordingfile.Close)
		}
	}
	return e
}

func (em *httpEmitter) WaitForProcessors() {
	em.wg.Wait()
}

func (em *httpEmitter) AddProcessor(p elo.Processor) {
	em.procs = append(em.procs, p)
	p.AddWaitGroup(em.wg)
	//go p.Loop()
}

func (em *httpEmitter) GetProcessor() []elo.Processor {
	return em.procs
}

func (em *httpEmitter) AddFilter(f elo.Filter) {
	em.filters = append(em.filters, f)
}

func (em *httpEmitter) GetFilters() []elo.Filter {
	return em.filters
}

func (em *httpEmitter) Loop() {
	const protocol = "http"
	port := em.config.Elo.Port
	handler := NewCsLogHandler(em, em.wg)
	go handler.Loop()

	//Build the address

	s := &http.Server{
		Addr:              ":" + port,
		Handler:           handler,
		ReadTimeout:       500 * time.Millisecond,
		ReadHeaderTimeout: 100 * time.Millisecond,
		WriteTimeout:      100 * time.Millisecond,
		IdleTimeout:       15 * time.Second,
		//MaxHeaderBytes: 1 << 20,

	}
	s.SetKeepAlivesEnabled(true)

	log.Info("Starting to listen on port %s", port)
	err := s.ListenAndServe()
	log.Info("Ended listening on port %s", port)

	if err != nil {
		log.Error("Error opening http server: %s", err)
	}
}

type csLogHandler struct {
	incoming chan transport
	emitter  *httpEmitter
	wg       *sync.WaitGroup
}

type transport struct {
	address string
	line    string
}

func NewCsLogHandler(em *httpEmitter, waitgroup *sync.WaitGroup) *csLogHandler {
	h := new(csLogHandler)
	h.emitter = em
	h.incoming = make(chan transport, em.config.Elo.BufferSize)
	h.wg = waitgroup
	return h
}

func (h *csLogHandler) receive(addr, line string) {
	if line == "" {
		log.Debug("csLogHandler received empty line.")
	} else {
		h.incoming <- transport{address: addr, line: line}
	}
}

func (h *csLogHandler) Loop() {
	log.Info("Starting csLogHandler Loop.")
	h.wg.Add(1)
	defer h.wg.Done()
	for {
		e := <-h.incoming
		h.pushMessage(e.address, e.line)
		if e.line == "cselo:StopRecorder" {
			break
		}
	}
	defer log.Info("Finishing csLogHandler Loop")
}

func (h *csLogHandler) pushMessage(remoteAddr, sbuf string) {
	if h.emitter.config.Elo.RecorderFileName != "" {
		h.emitter.recorder.Record(strings.Clone(sbuf))
	}
	t, m, err := elo.ShortenMessage(sbuf)
	if err != nil {
		log.Warn("Ignoring line. Error: %v.  Line: %s", err, sbuf)
		return
	} else {
		if !elo.CheckFilter(h.emitter, m) {
			for _, p := range h.emitter.procs {
				server := p.GetServer(remoteAddr)
				p.AddJob(elo.NewBaseEvent(server, t, m))
			}
		}
	}
}

func (h *csLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	var err error

	//log.Error("Got request: %s  Emitter: %p-%v  Server: %p-%v  Match: %v", r.URL, h.emitter, h.emitter, h.server, h.server, h.server.CurrentMatch)
	if r.Body == nil {
		log.Info("Empty request body. Url: %s", r.URL)
	} else {
		buf, err = io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Error("Problem reading request body: %v", err)
		} else {
			//remoteAddr := strings.Split(r.RemoteAddr, ":")[0]. // would be IPv4 address without port. Will not work with IPv6
			remoteAddr := "fromHttp" // should be IP of server=request IP. Unsolved problems with changing routes IPv4/IPv6, therefore only common sender
			h.receive(remoteAddr, string(buf))
		}
		w.WriteHeader(http.StatusOK)
	}
}
