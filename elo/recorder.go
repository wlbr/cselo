package elo

import (
	"bufio"
	"os"
	"sync"

	"github.com/wlbr/commons/log"
)

type Recorder struct {
	config   *Config
	fm       sync.Mutex //cm
	wg       *sync.WaitGroup
	incoming chan string
	filename string
	f        *os.File
	wbuf     *bufio.Writer
}

func NewRecorder(cfg *Config, waitgroup *sync.WaitGroup) *Recorder {
	//r := &Recorder{config: cfg, filename: cfg.Elo.RecorderFileName, cm: &sync.Mutex{}, fm: &sync.Mutex{}, incoming: make(chan string, cfg.Elo.BufferSize), wg: waitgroup}
	r := &Recorder{config: cfg, filename: cfg.Elo.RecorderFileName, incoming: make(chan string, cfg.Elo.BufferSize), wg: waitgroup}

	if cfg.Elo.RecorderFileName != "" {
		log.Info("Creating new Recorder with filename '%s'", r.filename)
		info, err := os.Stat(cfg.Elo.RecorderFileName)
		if err == nil {
			if info.IsDir() {
				log.Fatal("RecorderFileName is a directory, not overwriting.")
				cfg.FatalExit()
			}
			if !cfg.Elo.ForceOverwrite {
				//log.Fatal("RecorderFileName exists, not overwriting. Use -f to force overwrite.")
				//cfg.FatalExit()
				log.Warn("RecorderFileName exists, appending")
			}
		}

		f, err := os.OpenFile(r.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		if err != nil {
			log.Fatal("Could not create recorder file '%s': %s", cfg.Elo.RecorderFileName, err)
			cfg.FatalExit()
		} else {
			cfg.AddCleanUpFn(r.f.Close)
			r.f = f
			r.wbuf = bufio.NewWriter(r.f)
		}
	}
	return r
}

func (r *Recorder) writeLine(line string) {
	r.fm.Lock()
	defer r.fm.Unlock()
	log.Info("Recording: " + line)
	_, err := r.wbuf.WriteString(line)
	r.wbuf.Flush()
	if err != nil {
		log.Fatal("Could not write line to recorder file '%s': %s", r.filename, err)
		r.config.FatalExit()
	}
}

func (r *Recorder) Record(line string) {
	if r.config.Elo.RecorderFileName != "" {
		log.Info("Recorder getting job: %s", line)
		//r.cm.Lock()
		if len(line) > 0 && line[len(line)-1] != '\n' {
			r.incoming <- line + "\n"
		} else {
			r.incoming <- line
		}
		//r.cm.Unlock()
	}
}

func (r *Recorder) Loop() {
	if r.config.Elo.RecorderFileName == "" {
		log.Info("Skipping recorder loop, no filename given.")
	} else {
		log.Info("Starting recorder loop.")
		r.wg.Add(1)
		defer r.wg.Done()
		for {
			e := <-r.incoming
			r.writeLine(e)
			if e == "cselo:StopRecorder" {
				break
			}
		}
		defer log.Info("Finishing recorder")
	}
}
