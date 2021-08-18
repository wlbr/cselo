package postgresql

import (
	"context"
	"fmt"

	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"

	"github.com/jackc/pgx/v4"
)

type Postgres struct {
	config       *elo.Config
	Db           *pgx.Conn
	PlayersCache elo.PlayersCache
}

func NewPostgres(cfg *elo.Config) (*Postgres, error) {
	log.Info("Creating new PostgreSQL source")
	var err error
	s := &Postgres{config: cfg, PlayersCache: make(elo.PlayersCache)}

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
		s.Db, err = pgx.Connect(context.Background(), dbinfo)
		if err != nil {
			log.Error("Cannot open PostgresQL database: %v", err)
		}
		cfg.AddCleanUpFn(func() error {
			log.Info("Cleanup - closing PostgreSQL database connection")
			return s.Db.Close(context.Background())
		})
		log.Info("Established PostgreSQL database connection")
	} else {
		log.Warn("Not creating PostgreSQL source due to missing connection data")
		return nil, fmt.Errorf("Insufficient PostgreSQL connection data (%s)", err)
	}
	return s, err
}

/*
select initialname,count(actor) as kills,count(case when headshot=true then 1 end) as headshot, round(cast(count(case when headshot=true then 1 end) as float)/count(actor) * 1000)/10 as "hs%"  from kills
left join players on actor=players.id
WHERE timestmp > current_date - interval '30' day
group by initialname
order by count(actor) DESC;
*/
/*
func (s *Postgres) HandleKillEvent(e *events.Kill) {
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

func (s *Postgres) HandleAssistEvent(e *events.Assist) {
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
select initialname,count(case when victimtype='enemy' then 1 end) as enemyflashes, count(case when victimtype='teammate' then 1 end) as teammateflash, count(case when victimtype='self' then 1 end) as selfflash   from blindings
left join players on actor=players.id
WHERE timestmp > current_date - interval '30' day
group by initialname
order by count(case when victimtype='enemy' then 1 end) DESC
*/
/*
func (s *Postgres) HandleBlindedEvent(e *events.Blinded) {
	log.Info("Writing blind event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)
	object := s.GetOrStorePlayerbySteamID(e.Object)

	t := "enemy"
	switch {
	case e.SelfFlashed():
		t = "self"
		break
	case e.TeammateFlashed():
		t = "teammate"
		break
	}

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO blindings (actor, victim, duration,victimtype, timestmp) VALUES ($1, $2, $3, $4, $5)",
		subject.ID, object.ID, e.Duration, t, e.Time)
	if err != nil {
		log.Error("Cannot store BLINDING in PostgresQL database: %v", err)
	}
}

func (s *Postgres) HandleGrenadeEvent(e *events.Grenade) {
	log.Info("Writing grenade event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO grenadethrows (actor, grenadetype, timestmp) VALUES ($1, $2, $3)",
		subject.ID, e.GrenadeType, e.Time)
	if err != nil {
		log.Error("Cannot store GRENADETRHOW in PostgresQL database: %v", err)
	}
}

func (s *Postgres) HandlePlantedEvent(e *events.Planted) {
	log.Info("Writing planted event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO plantings (actor, timestmp) VALUES ($1, $2)",
		subject.ID, e.Time)
	if err != nil {
		log.Error("Cannot store PLANTING in PostgresQL database: %v", err)
	}
}

func (s *Postgres) HandleDefuseEvent(e *events.Defuse) {
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

func (s *Postgres) HandleBombedEvent(e *events.Bombed) {
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

func (s *Postgres) HandleHostageRescuedEvent(e *events.HostageRescued) {
	log.Info("Writing hostage rescued event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO rescues (actor, timestmp) VALUES ($1, $2)",
		subject.ID, e.Time)
	if err != nil {
		log.Error("Cannot store RESCUE in PostgresQL database: %v", err)
	}
}
func (s *Postgres) HandleRoundStartEvent(e *events.RoundStart) {}
func (s *Postgres) HandleRoundEndEvent(e *events.RoundEnd)     {}
func (s *Postgres) HandleMatchStartEvent(e *events.MatchStart) {}
func (s *Postgres) HandleMatchEndEvent(e *events.MatchEnd)     {}
*/
