package aggregators

import (
	"time"

	"github.com/wlbr/cs-elo/elo"
)

type Interval struct {
	start time.Time
	end   time.Time
}

type flashesPerPlayer struct {
	actor               *elo.Player
	teamflashes         int
	enemyflashes        int
	teamflashpercentage float32
}

type Flashes struct {
	timespan *Interval
	stats    []*flashesPerPlayer
}
