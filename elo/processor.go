package elo

type Processor interface {
	//AddWaitGroup(wg *sync.WaitGroup)
	AddJob(b *BaseEvent)
	Loop()
}
