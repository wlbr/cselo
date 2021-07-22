package sinks

import (
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
	"github.com/wlbr/cs-elo/elo/events"
)

type InfluxSink struct {
	config *elo.Config
	Client influxdb2.Client
	w      api.WriteAPI
	err    <-chan error
}

// // You can generate a Token from the "Tokens Tab" in the UI
// const influxdbtoken = "prh2UoJlsy02SWQQ2hhJTm_NGDo5uC_Q4QpY7DocWt7XyU-4FZSygKMrvSeW5NDWUU330CzIm1owZSYN353oGg=="

// const influxdbbucket = "cselo"
// const influxdborg = "kif"

func NewInfluxSink(cfg *elo.Config) (*InfluxSink, error) {
	is := &InfluxSink{config: cfg}
	err := false

	dbinfo := "http://"
	if cfg.InfluxDB.Host != "" {
		dbinfo += cfg.InfluxDB.Host
		if cfg.InfluxDB.Port != "" {
			dbinfo += ":" + cfg.InfluxDB.Port
		}
	} else {
		log.Warn("No InfluxDB host given")
		err = true
	}
	if cfg.InfluxDB.Token == "" {
		log.Warn("No InfluxDB token given")
	}
	is.Client = influxdb2.NewClient(dbinfo, cfg.InfluxDB.Token)

	if cfg.InfluxDB.Org == "" {
		log.Warn("No InfluxDB org given")
		err = true
	}
	if cfg.InfluxDB.Bucket == "" {
		log.Warn("No InfluxDB bucket given")
		err = true
	}
	if !err {
		is.w = is.Client.WriteAPI(cfg.InfluxDB.Org, cfg.InfluxDB.Bucket)
		cfg.AddCleanUpFn(func() error {
			is.w.Flush()
			log.Info("Cleanup - closing InfluxDB Connection")
			// always close client at the end
			is.Client.Close()
			return nil
		})

		is.err = is.w.Errors()
		go func() {
			for e := range is.err {
				log.Error("write error: %v\n", e)
			}
		}()
	} else {
		log.Warn("Not creating InfluxDB sink due to missing connection data")
		return nil, fmt.Errorf("Insufficient InfluxDB connection data")
	}
	return is, nil
}

func (s *InfluxSink) HandleKillEvent(e *events.Kill) {
	log.Info("Writing killevent to InfluxDB: %+v", e)
	h := 0
	if e.Headshot {
		h = 1
	}

	p := influxdb2.NewPointWithMeasurement("kills").
		AddTag("actor", e.Subject.Name).
		AddTag("victim", e.Object.Name).
		AddField("headshot", h).
		AddField("score", 1).
		SetTime(e.Time)
	// write point asynchronously
	s.w.WritePoint(p)

	// Flush writes
	//s.w.Flush()
}

// func (p *influxSink) HandleKillEvent2(e *KillEvent) {
// 	if strings.ToUpper(e.subject.SteamID) != "BOT" && strings.ToUpper(e.object.SteamID) != "BOT" {
// 		pkills := p.killStats[e.subject]
// 		if pkills == nil {
// 			pkills = AddPlayersKills(e.subject)
// 			p.killStats[e.subject] = pkills
// 		}
// 		pkills.Victims[e.object] = pkills.Victims[e.object] + 1
// 	}
// }

func (s *InfluxSink) HandleAssistEvent(e *events.Assist) {
	log.Info("Writing assist event to InfluxDB: %+v", e)

	p := influxdb2.NewPointWithMeasurement("assists").
		AddTag("actor", e.Subject.Name).
		AddTag("victim", e.Object.Name).
		AddField("score", 1).
		SetTime(e.Time)

	s.w.WritePoint(p)
}

func (s *InfluxSink) HandleBlindedEvent(e *events.Blinded) {
	log.Info("Writing blind event to InfluxDB: %+v", e)

	t := "enemy"
	switch {
	case e.SelfFlashed():
		t = "self"
		break
	case e.TeammateFlashed():
		t = "teammate"
		break
	}

	p := influxdb2.NewPointWithMeasurement("blinds").
		AddTag("actor", e.Subject.Name).
		AddTag("victim", e.Object.Name).
		AddTag("victimtype", t).
		AddField("score", 1).
		SetTime(e.Time)
	s.w.WritePoint(p)
}

func (s *InfluxSink) HandleGrenadeEvent(e *events.Grenade) {
	log.Info("Writing grenade event to InfluxDB: %+v", e)

	p := influxdb2.NewPointWithMeasurement("grenade").
		AddTag("actor", e.Subject.Name).
		AddTag("type", e.GrenadeType).
		AddField("score", 1).
		SetTime(e.Time)
	s.w.WritePoint(p)
}

func (s *InfluxSink) HandlePlantedEvent(e *events.Planted) {
	log.Info("Writing planted event to InfluxDB: %+v", e)

	p := influxdb2.NewPointWithMeasurement("planted").
		AddTag("actor", e.Subject.Name).
		AddField("score", 1).
		SetTime(e.Time)
	s.w.WritePoint(p)
}

func (s *InfluxSink) HandleDefuseEvent(e *events.Defuse) {
	log.Info("Writing defuse event to InfluxDB: %+v", e)

	p := influxdb2.NewPointWithMeasurement("defuse").
		AddTag("actor", e.Subject.Name).
		AddField("score", 1).
		SetTime(e.Time)
	s.w.WritePoint(p)
}

func (s *InfluxSink) HandleBombedEvent(e *events.Bombed) {
	log.Info("Writing bombed event to InfluxDB: %+v", e)

	p := influxdb2.NewPointWithMeasurement("bombed").
		AddTag("actor", e.Subject.Name).
		AddField("score", 1).
		SetTime(e.Time)
	s.w.WritePoint(p)
}

func (s *InfluxSink) HandleHostageRescuedEvent(e *events.HostageRescued) {
	log.Info("Writing hostage rescued event to InfluxDB: %+v", e)

	p := influxdb2.NewPointWithMeasurement("rescued").
		AddTag("actor", e.Subject.Name).
		AddField("score", 1).
		SetTime(e.Time)
	s.w.WritePoint(p)
}
func (s *InfluxSink) HandleRoundStartEvent(e *events.RoundStart) {}
func (s *InfluxSink) HandleRoundEndEvent(e *events.RoundEnd)     {}
func (s *InfluxSink) HandleMatchStartEvent(e *events.MatchStart) {}
func (s *InfluxSink) HandleMatchEndEvent(e *events.MatchEnd)     {}
func (s *InfluxSink) HandleGameOverEvent(e *events.GameOver)     {}
