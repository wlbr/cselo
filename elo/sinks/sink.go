package sinks

import (
	"github.com/wlbr/cselo/elo/events"
)

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
	HandleMatchStatusEvent(e *events.MatchStatus)
	HandleMatchEndEvent(e *events.MatchEnd)
	HandleAccoladeEvent(e *events.Accolade)
	HandleMatchCleanUpEvent(e *events.MatchCleanUp)
	HandleServerHibernationEvent(e *events.ServerHibernation)
	HandlePlayerConnectedEvent(e *events.PlayerConnected)
}
