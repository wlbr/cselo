package elo

import (
	"time"
)

type Processor interface {
	Dispatch(em Emitter, s string, t time.Time, m string)
}
