package main

import (
	"fmt"
	"time"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cselo/elo"
	"github.com/wlbr/cselo/elo/sources/postgresql"
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

	//aggregator := aggregators.
	// if s, e := sinks.NewInfluxSink(config); e == nil {
	// 	processor.AddSink(s)
	// }
	start := time.Now()

	if s, e := postgresql.NewPostgres(config); e == nil {
		ap := s.GetAllPlayers()
		for _, p := range ap {
			log.Info("player: %s", p)
		}

		p := s.PlayersCache[2]
		s.GetFlashStatsForPlayer(p, elo.IntervallLastWeek())
		s.GetFlashStatsForPlayer(p, elo.IntervallLastXDays(30))
		s.GetFlashStatsForPlayer(p, elo.IntervallLastXYears(1))
		s.GetFlashStatsForPlayer(p, elo.IntervallLastXYears(10))
		log.Warn("Player: %#v", p)

		//s.GetAllPlayersInInterval(time.Now().Add(time.Hour*24*-7), time.Now())

	}
	end := time.Now()

	elapsed := end.Sub(start)
	fmt.Printf("Processing took %s\n", elapsed)

	// for _, p := range emitter.GetProcessor() {
	// 	for player, ks := range p.GetKillStats() {
	// 		fmt.Printf("\nPLAYER %s killed:	 \n", player)
	// 		var kills []*playerkill
	// 		for victim, count := range ks.Victims {
	// 			kills = append(kills, &playerkill{player: victim.Name, count: count})
	// 		}
	// 		sort.Sort(ByCount(kills))
	// 		for _, p := range kills {
	// 			fmt.Printf("\t%s: %d \n", p.player, p.count)
	// 		}

	// 	}
	// }

	fmt.Printf("Processing took %s\n", elapsed)
}

/*
type playerkill struct {
	player string
	count  int
}

type ByCount []*playerkill

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Less(i, j int) bool { return a[i].count > a[j].count }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
*/
