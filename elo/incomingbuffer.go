package elo

import (
	"sync"

	"github.com/wlbr/commons/log"
)

type IncomingBuffer struct {
	config *Config
	im     *sync.Mutex
	// wg       *sync.WaitGroup
	incoming chan string
}

func NewIncomingBuffer(cfg *Config, waitgroup *sync.WaitGroup) *IncomingBuffer {
	b := &IncomingBuffer{config: cfg, incoming: make(chan string, cfg.Elo.BufferSize), im: &sync.Mutex{}} //, wg: waitgroup}, om: &sync.Mutex{},}
	//go b.Loop()
	return b
}

func (b *IncomingBuffer) Put(line string) {
	if b.config.Elo.RecorderFileName != "" {
		log.Info("Recorder getting job: %s", line)
		b.im.Lock()
		if len(line) > 0 && line[len(line)-1] != '\n' {
			b.incoming <- line + "\n"
		} else {
			b.incoming <- line
		}
		//b.incoming <- line
		b.im.Unlock()
	}
}

func (b *IncomingBuffer) Get() string {
	e := <-b.incoming
	return e
}

// func (b *IncomingBuffer) Loop() {
// 	log.Info("Starting incoming buffer loop.")
// 	b.wg.Add(1)
// 	defer b.wg.Done()
// 	for {
// 		e := <-b.incoming
// 		// b.writeLine(e) // add magic here
// 		if e == "cselo:StopProcessor" {
// 			break
// 		}
// 	}
// 	defer log.Info("Finishing incoming buffer")
// }
