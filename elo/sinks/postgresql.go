package sinks

import (
	"context"
	"fmt"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
	"github.com/wlbr/cs-elo/elo/events"

	"github.com/jackc/pgx/v4"
)

type PostgresSink struct {
	config *elo.Config
	db     *pgx.Conn
}

func NewPostgresSink(cfg *elo.Config) (*PostgresSink, error) {
	log.Info("Creating new PostgreSQL sink")
	var err error
	s := &PostgresSink{config: cfg}

	//pg connectstring "postgres://user:password@host:port5432/dbname"
	dbinfo := "postgres://"
	if cfg.PostgreSQL.User != "" {
		dbinfo += cfg.PostgreSQL.User
		if cfg.PostgreSQL.Password != "" {
			dbinfo += ":" + cfg.PostgreSQL.Password
		}
		dbinfo += "@"
	}
	if cfg.PostgreSQL.Host != "" {
		dbinfo += cfg.PostgreSQL.Host
		if cfg.PostgreSQL.Port != "" {
			dbinfo += ":" + cfg.PostgreSQL.Port
		}
	} else {
		log.Warn("No PostgresQL host given")
		err = fmt.Errorf("No PostgresQL host given")
	}
	if cfg.PostgreSQL.Database != "" {
		dbinfo += "/" + cfg.PostgreSQL.Database
	} else {
		log.Warn("No PostgresQL database name given")
		err = fmt.Errorf("No PostgresQL database name given")
	}
	if err == nil {
		s.db, err = pgx.Connect(context.Background(), dbinfo)
		if err != nil {
			log.Error("Cannot open PostgresQL database: %v", err)
		}
		cfg.AddCleanUpFn(func() error {
			log.Info("Cleanup - closing PostgreSQL database connection")
			return s.db.Close(context.Background())
		})
		log.Info("Established PostgreSQL database connection")
	} else {
		log.Warn("Not creating PostgreSQL sink due to missing connection data")
		return nil, fmt.Errorf("Insufficient PostgreSQL connection data (%s)", err)
	}
	return s, err
}

func (s *PostgresSink) GetPlayerIDBySteamID(steamid string) int64 {
	log.Info("Getting player fom PostgresDB. steamid='%v'", steamid)

	var id int64 = -1
	row := s.db.QueryRow(context.Background(), `SELECT id FROM players WHERE steamid=$1`,
		steamid)
	err := row.Scan(&id)
	switch {
	case err == pgx.ErrNoRows:
		log.Info("Cannot find PLAYER with steamid '%s' in PostgresQL database: %v", steamid, err)
	case err != nil:
		log.Error("Cannot read from PostgresQL database: %v", err)
	}
	return id
}

func (s *PostgresSink) GetOrStorePlayerbySteamID(p *elo.Player) *elo.Player {
	id := s.GetPlayerIDBySteamID(p.SteamID)
	if id == -1 {
		err := s.db.QueryRow(context.Background(), "INSERT INTO players  (initialname, steamid) VALUES ($1, $2) RETURNING id", p.Name, p.SteamID).Scan(&id)
		if err != nil {
			log.Error("Cannot store player in PostgresQL database: %v", err)
		}
		p.ID = id

	} else {
		p.ID = id
	}
	return p
}

/*
select initialname,count(actor) as kills,count(case when headshot=true then 1 end) as headshot, round(cast(count(case when headshot=true then 1 end) as float)/count(actor) * 1000)/10 as "hs%"  from kills
left join players on actor=players.id
WHERE timestmp > current_date - interval '30' day
group by initialname
order by count(actor) DESC;
*/
func (s *PostgresSink) HandleKillEvent(e *events.Kill) {
	log.Info("Writing killevent to PostgresDB: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)
	object := s.GetOrStorePlayerbySteamID(e.Object)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO kills (actor, victim, headshot, weapon, timestmp) VALUES ($1, $2, $3, $4, $5)",
		subject.ID, object.ID, e.Headshot, e.Weapon, e.Time)
	if err != nil {
		log.Error("Cannot store KILL in PostgresQL database: %v", err)
	}
}

func (s *PostgresSink) HandleAssistEvent(e *events.Assist) {
	log.Info("Writing assist event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)
	object := s.GetOrStorePlayerbySteamID(e.Object)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO assists (actor, victim, timestmp) VALUES ($1, $2, $3)",
		subject.ID, object.ID, e.Time)
	if err != nil {
		log.Error("Cannot store ASSIST in PostgresQL database: %v", err)
	}
}

/*
 select initialname,count(actor) as flashes, count(case when ownteam=true then 1 end) as teamflash, round(cast(count(case when ownteam=true then 1 end) as float)/count(actor) * 1000)/10 as "tf%"  from blindings
left join players on actor=players.id
WHERE timestmp > current_date - interval '30' day
group by initialname
order by round(cast(count(case when ownteam=true then 1 end) as float)/count(actor) * 1000)/10 DESC;
*/
func (s *PostgresSink) HandleBlindedEvent(e *events.Blinded) {
	log.Info("Writing blind event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)
	object := s.GetOrStorePlayerbySteamID(e.Object)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO blindings (actor, victim, duration,ownteam, timestmp) VALUES ($1, $2, $3, $4, $5)",
		subject.ID, object.ID, e.Duration, e.OwnTeam(), e.Time)
	if err != nil {
		log.Error("Cannot store BLINDING in PostgresQL database: %v", err)
	}
}

func (s *PostgresSink) HandleGrenadeEvent(e *events.Grenade) {
	log.Info("Writing grenade event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO grenadethrows (actor, grenadetype, timestmp) VALUES ($1, $2, $3)",
		subject.ID, e.GrenadeType, e.Time)
	if err != nil {
		log.Error("Cannot store GRENADETRHOW in PostgresQL database: %v", err)
	}
}

func (s *PostgresSink) HandlePlantedEvent(e *events.Planted) {
	log.Info("Writing planted event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO plantings (actor, timestmp) VALUES ($1, $2)",
		subject.ID, e.Time)
	if err != nil {
		log.Error("Cannot store PLANTING in PostgresQL database: %v", err)
	}
}

func (s *PostgresSink) HandleDefuseEvent(e *events.Defuse) {
	log.Info("Writing defuse event to PostgreSQL database: %+v", e)

	if e.Subject != nil {
		subject := s.GetOrStorePlayerbySteamID(e.Subject)
		_, err := s.db.Exec(context.Background(),
			"INSERT INTO defuses (actor, timestmp) VALUES ($1, $2)",
			subject.ID, e.Time)
		if err != nil {
			log.Error("Cannot store DEFUSE in PostgresQL database: %v", err)
		}
	}
}

func (s *PostgresSink) HandleBombedEvent(e *events.Bombed) {
	log.Info("Writing bombed event to PostgreSQL database: %+v", e)

	if e.Subject != nil {
		subject := s.GetOrStorePlayerbySteamID(e.Subject)
		_, err := s.db.Exec(context.Background(),
			"INSERT INTO bombings (actor, timestmp) VALUES ($1, $2)",
			subject.ID, e.Time)
		if err != nil {
			log.Error("Cannot store BOMBING in PostgresQL database: %v", err)
		}
	}
}

func (s *PostgresSink) HandleHostageRescuedEvent(e *events.HostageRescued) {
	log.Info("Writing hostage rescued event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO rescues (actor, timestmp) VALUES ($1, $2)",
		subject.ID, e.Time)
	if err != nil {
		log.Error("Cannot store RESCUE in PostgresQL database: %v", err)
	}
}
func (s *PostgresSink) HandleRoundStartEvent(e *events.RoundStart) {}
func (s *PostgresSink) HandleRoundEndEvent(e *events.RoundEnd)     {}
func (s *PostgresSink) HandleMatchStartEvent(e *events.MatchStart) {}
func (s *PostgresSink) HandleMatchEndEvent(e *events.MatchEnd)     {}
