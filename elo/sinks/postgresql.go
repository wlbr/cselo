package sinks

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/wlbr/commons/log"
	"github.com/wlbr/cs-elo/elo"
	"github.com/wlbr/cs-elo/elo/events"
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
		if p.ProfileID == "" {
			p.ProfileID = elo.SteamIdToProfileId(p.SteamID)
		}
		err := s.db.QueryRow(context.Background(), "INSERT INTO players  (initialname, steamid, profileid) VALUES ($1, $2, $3) RETURNING id", p.Name, p.SteamID, p.ProfileID).Scan(&id)
		if err != nil {
			log.Error("Cannot store player '%+v' in PostgresQL database: %v", p, err)
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
WHERE timestmp > current_date - interval '2' day
group by initialname
order by count(actor) DESC;
*/
func (s *PostgresSink) HandleKillEvent(e *events.Kill) {
	log.Info("Writing killevent to PostgresDB: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)
	object := s.GetOrStorePlayerbySteamID(e.Object)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO kills (match, actor, victim, headshot, weapon, timestmp) VALUES ($1, $2, $3, $4, $5, $6)",
		e.Server.CurrentMatch.ID, subject.ID, object.ID, e.Headshot, e.Weapon, e.Time)
	if err != nil {
		log.Error("Cannot store KILL in PostgresQL database: %v", err)
	}
}

func (s *PostgresSink) HandleAssistEvent(e *events.Assist) {
	log.Info("Writing assist event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)
	object := s.GetOrStorePlayerbySteamID(e.Object)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO assists (match, actor, victim, timestmp) VALUES ($1, $2, $3, $4)",
		e.Server.CurrentMatch.ID, subject.ID, object.ID, e.Time)
	if err != nil {
		log.Error("Cannot store ASSIST in PostgresQL database: %v", err)
	}
}

/*
select initialname,count(case when victimtype='enemy' then 1 end) as enemyflashes, count(case when victimtype='teammate' then 1 end) as teammateflash, count(case when victimtype='self' then 1 end) as selfflash   from blindings
left join players on actor=players.id
WHERE timestmp > current_date - interval '2' day
group by initialname
order by count(case when victimtype='enemy' then 1 end) DESC;
*/
func (s *PostgresSink) HandleBlindedEvent(e *events.Blinded) {
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
		"INSERT INTO blindings (match, actor, victim, duration,victimtype, timestmp) VALUES ($1, $2, $3, $4, $5, $6)",
		e.Server.CurrentMatch.ID, subject.ID, object.ID, e.Duration, t, e.Time)
	if err != nil {
		log.Error("Cannot store BLINDING in PostgresQL database: %v", err)
	}
}

/*
SELECT initialname,count(case when grenadetype='flashbang' then 1 end) as flash,
    count(case when grenadetype='hegrenade' then 1 end) as he,
	count(case when grenadetype='molotov' then 1 end) as molotov,
	count(case when grenadetype='smoke' then 1 end) as smoke,
	count(case when grenadetype='decoy' then 1 end) as decoy FROM grenadethrows
LEFT JOIN players ON actor=players.id
WHERE timestmp > current_date - interval '3' day
GROUP BY initialname
ORDER BY flash DESC;
*/

func (s *PostgresSink) HandleGrenadeEvent(e *events.Grenade) {
	log.Info("Writing grenade event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO grenadethrows (match, actor, grenadetype, timestmp) VALUES ($1, $2, $3, $4)",
		e.Server.CurrentMatch.ID, subject.ID, e.GrenadeType, e.Time)
	if err != nil {
		log.Error("Cannot store GRENADETRHOW in PostgresQL database: %v", err)
	}
}

func (s *PostgresSink) HandlePlantedEvent(e *events.Planted) {
	log.Info("Writing planted event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO scoreaction (match, actor, actiontype, timestmp) VALUES ($1, $2, $3, $4)",
		e.Server.CurrentMatch.ID, subject.ID, "planting", e.Time)
	if err != nil {
		log.Error("Cannot store PLANTING in PostgresQL database: %v", err)
	}
}

func (s *PostgresSink) HandleDefuseEvent(e *events.Defuse) {
	log.Info("Writing defuse event to PostgreSQL database: %+v", e)

	if e.Subject != nil {
		subject := s.GetOrStorePlayerbySteamID(e.Subject)
		_, err := s.db.Exec(context.Background(),
			"INSERT INTO scoreaction (match, actor, actiontype, timestmp) VALUES ($1, $2, $3, $4)",
			e.Server.CurrentMatch.ID, subject.ID, "defuse", e.Time)
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
			"INSERT INTO scoreaction (match, actor, actiontype, timestmp) VALUES ($1, $2, $3, $4)",
			e.Server.CurrentMatch.ID, subject.ID, "bombing", e.Time)
		if err != nil {
			log.Error("Cannot store BOMBING in PostgresQL database: %v", err)
		}
	}
}

func (s *PostgresSink) HandleHostageRescuedEvent(e *events.HostageRescued) {
	log.Info("Writing hostage rescued event to PostgreSQL database: %+v", e)

	subject := s.GetOrStorePlayerbySteamID(e.Subject)

	_, err := s.db.Exec(context.Background(),
		"INSERT INTO scoreaction (match, actor, actiontype, timestmp) VALUES ($1, $2, $3, $4)",
		e.Server.CurrentMatch.ID, subject.ID, "rescue", e.Time)
	if err != nil {
		log.Error("Cannot store RESCUE in PostgresQL database: %v", err)
	}
}

// "INSERT INTO matches (gamemode, mapgroup, mapfullname, mapname, scorea, scoreb, duration, matchend, timestmp) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
// e.GameMode, e.MapGroup, e.MapFullName, e.MapName, e.ScoreA, e.ScoreB, e.Duration, e.MatchEnd, e.Time)
func (s *PostgresSink) HandleMatchEndEvent(e *events.MatchEnd) {
	log.Info("Writing game over event to PostgreSQL database: %+v", e)
	m := e.Server.CurrentMatch
	_, err := s.db.Exec(context.Background(), `UPDATE matches
		SET gamemode=$2, mapgroup=$3, mapfullname=$4, mapname=$5, scorea=$6, scoreb=$7, duration=$8, matchend=$9, timestmp=$10
		WHERE id=$1`,
		m.ID, m.GameMode, m.MapGroup, m.MapFullName, m.MapName, m.ScoreA, m.ScoreB, m.Duration, m.End, e.Time)
	if err != nil {
		log.Error("Cannot store MATCHEND in PostgresQL database: %v  message:`%s'", err, e.Message)
	}
}

func (s *PostgresSink) HandleMatchStartEvent(e *events.MatchStart) {
	log.Info("Writing game start event to PostgreSQL database: %+v", e)
	var id int64
	err := s.db.QueryRow(context.Background(),
		`INSERT INTO matches (mapfullname, mapname, matchstart, timestmp) VALUES ($1, $2, $3, $4)
		RETURNING id`,
		e.MapFullName, e.MapName, e.Time, e.Time).Scan(&id)
	e.Server.CurrentMatch.ID = id
	if err != nil {
		log.Error("Cannot store MATCHSTART in PostgresQL database: %v  message:`%s'", err, e.Message)
	}
}

func (s *PostgresSink) HandleRoundStartEvent(e *events.RoundStart) {}
func (s *PostgresSink) HandleRoundEndEvent(e *events.RoundEnd)     {}

func (s *PostgresSink) HandleAccoladeEvent(e *events.Accolade) {
	log.Info("Writing accolade event to PostgreSQL database: %+v", e)
	//fmt.Printf("Accolade: %s - %s\n", e.Subject.Name, e.Type)
}
