package elo

import "sync"

type Processor interface {
	AddWaitGroup(wg *sync.WaitGroup)
	AddJob(b *BaseEvent)
	Loop()
	GetServer(ip string) *Server
}
