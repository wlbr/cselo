package main

import (
	"fmt"
	"time"

	"github.com/wlbr/cs-elo/elo"
	"github.com/wlbr/cs-elo/elo/processors"
	"github.com/wlbr/cs-elo/elo/sinks"
)

var (
	//Version is a linker injected variable for a git revision info used as version info
	Version = "Unknown build"
	//BuildTimestamp is a linker injected variable for a buildtime timestamp used in version info
	BuildTimestamp = "unknown build timestamp."

	config *elo.Config
)

func main() {

	config = new(elo.Config)
	config.Initialize(Version, BuildTimestamp)
	defer config.CleanUp()

	processor := processors.NewCsgoLogProcessor(config)
	// if s, e := sinks.NewInfluxSink(config); e == nil {
	// 	processor.AddSink(s)
	// }
	if s, e := sinks.NewPostgresSink(config); e == nil {
		processor.AddSink(s)
	}
	// if s, e := sinks.NewPrinterSink(config); e == nil {
	// 	processor.AddSink(s)
	// }

	var emitter elo.Emitter
	if config.Elo.CsLogFileName != "" {
		emitter = elo.NewFileEmitter(config)
	} else {
		emitter = elo.NewUdpEmitter(config)
	}
	emitter.AddFilter(&elo.AllBotsFilter{})
	emitter.AddProcessor(processor)

	start := time.Now()

	emitter.Loop()

	end := time.Now()
	elapsed := end.Sub(start)
	fmt.Printf("Processing took %s\n", elapsed)

	fmt.Printf("Processing took %s\n", elapsed)
}

type playerkill struct {
	player string
	count  int
}

type ByCount []*playerkill

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Less(i, j int) bool { return a[i].count > a[j].count }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
