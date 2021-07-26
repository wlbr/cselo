package elo

import (
	"time"
)

type Processor interface {
	Dispatch(em Emitter, s *Server, t time.Time, m string)
}
