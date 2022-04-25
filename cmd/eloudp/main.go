package main

import (
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/processors"
	"github.com/wlbr/cselo/elo/sinks"
	"github.com/wlbr/cselo/net"
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
	log.Warn("Starting up")

	processor := processors.NewCsgoLogProcessor(config)
	// if s, e := sinks.NewInfluxSink(config); e == nil {
	// 	processor.AddSink(s)
	// }

	discord := net.NewDisordSender(config.Elo.DiscordWebhook)
	if config.Elo.ImportFileName != "" {
		discord = net.NewDisordSender("") //inactive
	}

	if s, e := sinks.NewPostgresSink(config, discord); e == nil {
		processor.AddSink(s)
	}
	// if s, e := sinks.NewPrinterSink(config); e == nil {
	// 	processor.AddSink(s)
	// }

	var emitter elo.Emitter
	if config.Elo.ImportFileName != "" {
		emitter = elo.NewFileEmitter(config)
	} else {
		emitter = elo.NewUdpEmitter(config)
	}
	emitter.AddFilter(&elo.AllBotsFilter{})
	emitter.AddFilter(&elo.UnknownFilter{})
	emitter.AddProcessor(processor)

	start := time.Now()

	emitter.Loop()

	end := time.Now()
	elapsed := end.Sub(start)

	defer log.Warn("Shutting down, been up for %s\n", elapsed)
}

type playerkill struct {
	player string
	count  int
}

type ByCount []*playerkill

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Less(i, j int) bool { return a[i].count > a[j].count }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
