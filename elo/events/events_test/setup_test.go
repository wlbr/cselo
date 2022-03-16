package events

import (
	"flag"
	"os"
	"testing"
	"time"

	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/processors"
	"github.com/wlbr/cselo/elo/sinks"
	"github.com/wlbr/cselo/elo/sources/postgresql"
)

var config *elo.Config
var store *postgresql.Postgres
var testfile string
var player string

var counter *sinks.InMemoryCounterSink

func TestMain(m *testing.M) {
	testsetup()
	testsetupDB()
	code := m.Run()
	testteardown()
	os.Exit(code)
}

func testsetup() {
	flag.StringVar(&player, "player", "Jagger", "Player to try the tests on.")

	config = new(elo.Config)
	config.Initialize("Test", time.Now().Format(time.ANSIC))

	testfile = config.Elo.CsLogFileName

	processor := processors.NewCsgoLogProcessor(config)
	if s, e := sinks.NewInMemoryCounterSink(config, player); e == nil {
		processor.AddSink(s)
		counter = s
	}

	emitter := elo.NewFileEmitter(config)

	emitter.AddFilter(&elo.AllBotsFilter{})
	emitter.AddFilter(&elo.UnknownFilter{})
	emitter.AddProcessor(processor)

	emitter.Loop()

}

func testsetupDB() {

	// db tests
	var err error
	store, err = postgresql.NewPostgres(config)
	if err != nil {
		panic(err.Error())
	}

}

func testteardown() {
	defer config.CleanUp()
}
