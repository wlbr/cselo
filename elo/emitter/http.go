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
	//wbuf          *bufio.Writer
	filters  []elo.Filter
	recorder *elo.Recorder
	incoming *elo.IncomingBuffer
	server   *elo.Server
}

func NewHttpEmitter(cfg *elo.Config) *httpEmitter {
	e := new(httpEmitter)
	e.wg = new(sync.WaitGroup)
	//e.m = new(sync.Mutex)
	// e.incoming = elo.NewIncomingBuffer(cfg, e.wg)
	// go e.pusher()
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
			//e.wbuf = bufio.NewWriter(e.recordingfile)
		}
	}
	return e
}

func (em *httpEmitter) WaitForProcessors() {
	em.wg.Wait()
}

func (em *httpEmitter) AddProcessor(p elo.Processor) {
	em.procs = append(em.procs, p)
	//p.AddWaitGroup(em.wg)
	go p.Loop()
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

func (em *httpEmitter) stopWorkers(server *elo.Server) {
	for _, p := range em.procs {
		p.AddJob(elo.NewBaseEvent(server, time.Now(), "cselo:StopProcessing."))
	}

	//wait for the processors to stop
	em.WaitForProcessors()
}

// func (em *httpEmitter) pusher() {
// 	em.wg.Add(1)
// 	defer em.wg.Done()
// 	for {
// 		sbuf := em.incoming.Get()

// 		t, m, err := elo.ShortenMessage(sbuf)
// 		if err != nil {
// 			log.Warn("Ignoring line. %v", err)
// 			return
// 		} else {
// 			if !elo.CheckFilter(em, m) {
// 				for _, p := range em.procs {
// 					p.AddJob(elo.NewBaseEvent(em.server, t, m))
// 				}
// 			} else {
// 				log.Info("Filtered Message: %s", sbuf)
// 			}
// 		}
// 		if sbuf == "cselo:StopProcessor" {
// 			break
// 		}
// 	}
// }

func (em *httpEmitter) Loop() {
	const protocol = "http"
	port := em.config.Elo.Port
	server := elo.NewServer("fromHttp")
	em.server = server
	defer em.stopWorkers(server)
	handler := &csLogHandler{emitter: em, server: server}

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
	// go func() {
	// 	if err := http.ListenAndServe(":"+port, handler.HandleFastHTTP); err != nil {
	// 		log.Fatal("error in ListenAndServe: %v", err)
	// 	}
	// }()
	// srv := fasthttp.Server{
	// 	Concurrency: 10, // Number of concurrent connections ,
	// }

	// srv.Handler = handler.HandleFastHTTP
	// //err := srv.ListenAndServe(":" + port)
	// ln, err := net.Listen("tcp4", ":"+port)
	// err = fasthttp.Serve(ln, handler.HandleFastHTTP)

	if err != nil {
		log.Error("Error opening http server: %s", err)
	}
}

type csLogHandler struct {
	emitter *httpEmitter
	server  *elo.Server
}

// // request handler in net/http style, i.e. method bound to csLogHandler struct.
// func (h *csLogHandler) HandleFastHTTP(ctx *fasthttp.RequestCtx) {
// 	// notice that we may access MyHandler properties here - see h.foobar.
// 	var buf []byte
// 	var err error

// 	buf = ctx.PostBody()
// 	//log.Error("Got request: %s  Emitter: %p-%v  Server: %p-%v  Match: %v", r.URL, h.emitter, h.emitter, h.server, h.server, h.server.CurrentMatch)
// 	if buf == nil {
// 		log.Info("Empty request body. Url: %s", ctx.URI())
// 		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
// 	} else {
// 		ctx.SetStatusCode(fasthttp.StatusOK)

// 		//fmt.Println(string(buf))
// 		if err != nil {
// 			log.Error("Problem reading request body: %v", err)
// 		} else {
// 			sbuf := string(buf)
// 			go h.pushMessage(sbuf)
// 		}
// 	}
// 	fmt.Fprintf(ctx, "")
// }

func (h *csLogHandler) pushMessage(sbuf string) {
	t, m, err := elo.ShortenMessage(sbuf)
	if err != nil {
		log.Warn("Ignoring line. %v", err)
		return
	} else {
		if !elo.CheckFilter(h.emitter, m) {
			for _, p := range h.emitter.procs {
				p.AddJob(elo.NewBaseEvent(h.server, t, m))
			}
		}
		if h.emitter.config.Elo.RecorderFileName != "" {
			h.emitter.recorder.Record(sbuf)
		}
	}
}

func (h *csLogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var buf []byte
	var err error

	//defer r.Body.Close()
	//log.Error("Got request: %s  Emitter: %p-%v  Server: %p-%v  Match: %v", r.URL, h.emitter, h.emitter, h.server, h.server, h.server.CurrentMatch)
	if r.Body == nil {
		log.Info("Empty request body. Url: %s", r.URL)
	} else {
		buf, err = io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Error("Problem reading request body: %v", err)
		} else {
			h.pushMessage(strings.Clone(string(buf)))
		}

	}
}
