package aggregators

import (
	"time"

	"github.com/wlbr/cselo/elo"
)

type Interval struct {
	start time.Time
	end   time.Time
}

type flashesPerPlayer struct {
	actor             *elo.Player
	enemyflashes      int
	teammateflashes   int
	selfflashed       int
	enemyflashPrct    float32
	teammateflashPrct float32
	selfflashPrct     float32
	throws            int
}

type Flashes struct {
	timespan *Interval
	stats    []*flashesPerPlayer
}
