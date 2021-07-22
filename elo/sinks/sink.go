package sinks

import "github.com/wlbr/cs-elo/elo/events"

type Sink interface {
	HandleKillEvent(e *events.Kill)
	HandleAssistEvent(e *events.Assist)
	HandleBlindedEvent(e *events.Blinded)
	HandleGrenadeEvent(e *events.Grenade)
	HandlePlantedEvent(e *events.Planted)
	HandleDefuseEvent(e *events.Defuse)
	HandleBombedEvent(e *events.Bombed)
	HandleHostageRescuedEvent(e *events.HostageRescued)
	HandleRoundStartEvent(e *events.RoundStart)
	HandleRoundEndEvent(e *events.RoundEnd)
	HandleMatchStartEvent(e *events.MatchStart)
	HandleMatchEndEvent(e *events.MatchEnd)
	HandleGameOverEvent(e *events.GameOver)
}
